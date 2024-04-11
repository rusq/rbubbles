package filemgr

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bbtea/display"
)

type Model struct {
	Globs     []string
	Selected  string
	FS        fs.FS
	Directory string
	Height    int
	ShowHelp  bool
	Style     Style
	files     []fs.FileInfo
	finished  bool
	st        display.State
	viewStack display.Stack[display.State]

	Debug bool
	last  string // last key pressed
}

type Style struct {
	Normal    lipgloss.Style
	Directory lipgloss.Style
	Inverted  lipgloss.Style
}

func New(fsys fs.FS, dir string, height int, globs ...string) Model {
	return Model{
		Globs:     globs,
		FS:        fsys,
		Directory: dir,
		Height:    height,
		Style: Style{
			Normal:    lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
			Directory: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
			Inverted:  lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Background(lipgloss.Color("240")),
		},
	}
}

type wmReadDir struct {
	dir   string
	files []fs.FileInfo
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return readFS(m.FS, m.Directory, m.Globs...)
	}
}

func readFS(fsys fs.FS, dir string, globs ...string) wmReadDir {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return wmReadDir{dir, nil}
	}
	dirs, err := collectDirs(sub)
	if err != nil {
		return wmReadDir{dir, nil}
	}
	files, err := collectFiles(sub, globs...)
	if err != nil {
		return wmReadDir{dir, nil}
	}
	if !(dir == "." || dir == "/" || dir == "") {
		files = append([]fs.FileInfo{specialDir{".."}}, files...)
	}
	return wmReadDir{dir, append(files, dirs...)}
}

func collectFiles(fsys fs.FS, globs ...string) (files []fs.FileInfo, err error) {
	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			// do not show the current directory
			return nil
		}
		if d.IsDir() {
			return fs.SkipDir
		}
		for _, glob := range globs {
			if ok, err := filepath.Match(glob, d.Name()); err != nil {
				return err
			} else if ok {
				fi, err := d.Info()
				if err != nil {
					return err
				}
				files = append(files, fi)
			}
		}
		return nil
	})
	return
}

func collectDirs(fsys fs.FS) ([]fs.FileInfo, error) {
	var dirs []fs.FileInfo
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
		if d.IsDir() {
			dir, err := d.Info()
			if err != nil {
				return err
			}
			dirs = append(dirs, dir)
			return fs.SkipDir
		}
		return nil
	})
	return dirs, err
}

func (m Model) height() int {
	if m.ShowHelp {
		return m.Height - 2
	}
	return m.Height
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case error:
		log.Printf("error: %v", msg)
		return m, tea.Quit
	case tea.WindowSizeMsg:
		if m.Height == 0 {
			m.Height = msg.Height
		}
	case tea.KeyMsg:
		m.last = msg.String()
		switch msg.String() {
		case "ctrl+c", "q":
			m.finished = true
			return m, tea.Quit
		case "up", "ctrl+p", "k":
			m.st.Up()
		case "down", "ctrl+n", "j":
			m.st.Down(len(m.files))
		case "right", "pgdown", "ctrl+v", "ctrl+f":
			m.st.NextPg(m.height(), len(m.files))
		case "left", "pgup", "alt+v", "ctrl+b":
			m.st.PrevPg(m.height())
		case "ctrl+r":
			return m, tea.Batch(m.Init())
		case "enter", "ctrl+m":
			if len(m.files) == 0 {
				break
			}
			if m.files[m.st.Cursor].IsDir() {
				m.Directory = filepath.Join(m.Directory, m.files[m.st.Cursor].Name())
				m.viewStack.Push(m.st)
				m.st = display.State{}
				return m, tea.Batch(m.Init())
			}
			cmds = append(cmds, selectedCmd(m.Directory, m.files[m.st.Cursor]))
		case "backspace", "ctrl+h":
			if m.viewStack.Len() > 0 {
				m.st = m.viewStack.Pop()
				m.Directory = filepath.Dir(m.Directory)
				return m, tea.Batch(m.Init())
			}
		}
	case wmReadDir:
		m.files = msg.files
		m.st.SetMax(m.height())
	}

	return m, tea.Batch(cmds...)
}

func selectedCmd(dir string, fi fs.FileInfo) tea.Cmd {
	return func() tea.Msg {
		return WMSelected{
			Filepath: filepath.Join(dir, fi.Name()),
			IsDir:    fi.IsDir(),
		}
	}
}

type WMSelected struct {
	Filepath string
	IsDir    bool
}

// humanizeSize returns a human-readable string representing a file size.
// for example 240.4M or 2.3G
func humanizeSize(size int64) string {
	const (
		K = 1 << 10
		M = 1 << 20
		G = 1 << 30
		T = 1 << 40
	)

	switch {
	case size < K:
		return fmt.Sprintf("%5dB", size)
	case size < M:
		return fmt.Sprintf("%5.1fK", float64(size)/K)
	case size < G:
		return fmt.Sprintf("%5.1fM", float64(size)/M)
	case size < T:
		return fmt.Sprintf("%5.1fG", float64(size)/G)
	default:
		return fmt.Sprintf("%5.1fT", float64(size)/G)

	}
}

const Width = 40

func printFile(fi fs.FileInfo) string {
	// filename.extension  <DIR>  02-01-2006 15:04
	const (
		dttmLayout = "02-01-2006 15:04"
		dirMarker  = "<DIR>"
		filesizeSz = 6
		dttmSz     = len(dttmLayout)
		filenameSz = Width - filesizeSz - dttmSz - 3
	)

	var sz = dirMarker
	if !fi.IsDir() {
		sz = humanizeSize(fi.Size())
	}
	return fmt.Sprintf("%-*s %*s %s", filenameSz, display.Trunc(fi.Name(), filenameSz), filesizeSz, sz, fi.ModTime().Format(dttmLayout))
}

func (m Model) printDebug(w io.Writer) {
	fmt.Fprintf(w, "cursor: %d\n", m.st.Cursor)
	fmt.Fprintf(w, "min: %d\n", m.st.Min)
	fmt.Fprintf(w, "max: %d\n", m.st.Max)
	fmt.Fprintf(w, "last: %q\n", m.last)
	fmt.Fprintf(w, "dir: %q\n", m.Directory)
	fmt.Fprintf(w, "selected: %q\n", m.Selected)
	for i := range Width {
		if i%10 == 0 {
			w.Write([]byte{'|'})
		} else {
			fmt.Fprint(w, i%10)
		}
	}
	fmt.Fprintln(w)
}

func (m Model) View() string {
	if m.finished {
		return ""
	}
	var buf strings.Builder
	if m.Debug {
		m.printDebug(&buf)
	}
	if len(m.files) == 0 {
		buf.WriteString(m.Style.Normal.Render("No files found, press [Backspace]") + "\n")
		for i := 0; i < m.height()-1; i++ {
			fmt.Fprintln(&buf, m.Style.Normal.Render(strings.Repeat(" ", Width-1))) //padding
		}
	} else {
		for i, file := range m.files {
			if i < m.st.Min || i > m.st.Max {
				continue
			}
			style := m.Style.Normal
			if file.IsDir() {
				style = m.Style.Directory
			}
			if i == m.st.Cursor {
				style = m.Style.Inverted
			}
			fmt.Fprintln(&buf, style.Render(printFile(file)))
		}
		numDisplayed := m.st.Displayed(len(m.files))
		for i := 0; i < m.height()-numDisplayed; i++ {
			fmt.Fprintln(&buf, m.Style.Normal.Render(strings.Repeat(" ", Width-1)))
		}
	}
	if m.ShowHelp {
		buf.WriteString("\n ↑ ↓ move・[⏎] select・[⇤] back・[q] quit\n")
	}
	return buf.String()
}

type specialDir struct {
	name string
}

func (s specialDir) Name() string {
	return s.name
}

func (s specialDir) Size() int64 {
	return 0
}

func (s specialDir) Mode() fs.FileMode {
	return fs.ModeDir
}

func (s specialDir) ModTime() time.Time {
	return time.Time{}
}

func (s specialDir) IsDir() bool {
	return true
}

func (s specialDir) Sys() interface{} {
	return s
}
