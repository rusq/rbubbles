package customise

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"bbtea/display"
	"bbtea/filemgr"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Items     []Item
	nameColSz int
	Cursor    string
	width     int
	editing   bool
	edittype  Type
	finishing bool
	st        display.State
	err       error
	Style     Styles
	fields
}

type Styles struct {
	Normal      lipgloss.Style
	Selected    lipgloss.Style
	Description lipgloss.Style
}

type fields struct {
	textarea  textarea.Model
	textinput textinput.Model
	radio     RadioButton
	filemgr   filemgr.Model
}

var (
	defStyle = lipgloss.NewStyle()
)

func NewModel(items []Item) Model {
	maxNameLen := 0
	for i := range items {
		if l := len(items[i].Name()); maxNameLen < l {
			maxNameLen = l
		}
	}
	return Model{
		Items:     items,
		nameColSz: maxNameLen,
		Cursor:    "",
		Style: Styles{
			Normal:      defStyle,
			Selected:    defStyle.Copy().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("7")),
			Description: defStyle.Copy().Faint(true),
		},
		fields: fields{
			textarea:  textarea.New(),
			textinput: textinput.New(),
			radio:     RadioButton{},
			filemgr:   filemgr.New(os.DirFS("."), ".", 0, "*"),
		},
	}
}

func (m Model) Init() tea.Cmd {
	slog.Debug("customise.init")
	return tea.Batch(m.filemgr.Init())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.editing {
		return m.procMsgEdit(msg)
	}
	return m.procMsgView(msg)
}

func (m Model) procMsgView(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.st.SetMax(msg.Height)
		m.filemgr.Height = 10
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.st.Down(len(m.Items))
		case "k", "up":
			m.st.Up()
		case "home":
			m.st.Home(len(m.Items))
		case "end":
			m.st.End(len(m.Items), len(m.Items))
		case " ":
			if m.Items[m.st.Cursor].Type() != TCheckbox {
				break
			}
			fallthrough
		case "enter", "f4":
			item := m.Items[m.st.Cursor]

			m.edittype = item.Type()
			switch m.edittype {
			case TMultiline:
				m.textarea.Reset()
				m.textarea.SetValue(item.Value())
				m.textarea.Focus()
				m.editing = true
			case TText:
				m.textinput.Reset()
				m.textinput.SetValue(item.Value())
				m.textinput.Focus()
				m.editing = true
			case TRadio:
				m.radio.SetValues(item.AllowedValues(), item.Value())
				m.editing = true
			case TFileExisting:
				m.editing = true
				m.filemgr.Focus()
				cmds = append(cmds, m.filemgr.Init())
				m.filemgr.Select(item.Value())
			case TCheckbox:
				if item.Value() == sTrue {
					item.Set(sFalse)
				} else {
					item.Set(sTrue)
				}
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) procMsgEdit(msg tea.Msg) (Model, tea.Cmd) {
OUTER:
	switch msg := msg.(type) {
	case filemgr.WMSelected:
		if !msg.IsDir {
			m.Items[m.st.Cursor].Set(msg.Filepath)
			m.editing = false
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// we only process enter for non-multiline modes
			switch m.edittype {
			case TMultiline, TFileExisting:
				break OUTER
			}
			fallthrough
		case "esc":
			m.editing = false
			var val string
			switch m.edittype {
			case TText:
				m.textinput.Blur()
				val = m.textinput.Value()
			case TMultiline:
				m.textarea.Blur()
				val = m.textarea.Value()
			case TRadio:
				val = m.radio.Value()
			case TFileExisting:
				if m.filemgr.Selected == "" {
					val = m.Items[m.st.Cursor].Value()
				} else {
					val = m.filemgr.Selected
				}
				m.filemgr.Blur()
			}
			m.Items[m.st.Cursor].Set(val)
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch m.edittype {
	case TText:
		m.textinput, cmd = m.textinput.Update(msg)
	case TMultiline:
		m.textarea, cmd = m.textarea.Update(msg)
	case TRadio:
		m.radio, cmd = m.radio.Update(msg)
	case TFileExisting:
		m.filemgr, cmd = m.filemgr.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.finishing {
		return ""
	}
	if m.err != nil {
		return m.err.Error()
	}
	if m.editing && m.edittype != TCheckbox {
		return m.editView()
	} else {
		return m.selectView()
	}
}

func (m Model) selectView() string {
	if len(m.Items) == 0 {
		return "No items to show."
	}
	var (
		empty = strings.Repeat(" ", len(m.Cursor))
	)

	var buf strings.Builder

	for i, item := range m.Items {
		cursor := empty
		if m.st.IsSelected(i) {
			cursor = m.Cursor
		}
		value := item.Value()
		if len(value) == 0 {
			value = "<empty>"
		}
		var val string
		switch item.Type() {
		case TMultiline, TText, TFileExisting:
			val = display.Trunc(value, m.width)
		case TCheckbox:
			if value == sTrue {
				val = "[x]"
			} else {
				val = "[ ]"
			}
		case TRadio:
			val = "[" + display.Trunc(value, m.width-4) + " ↓]"
		}
		style := m.Style.Normal
		if m.st.IsSelected(i) {
			style = m.Style.Selected
		}
		fmt.Fprintf(&buf,
			style.Render("%s%*s  %v")+"\n",
			cursor,
			-m.nameColSz,
			item.Name(),
			val,
		)
	}

	// description
	fmt.Fprint(&buf, "\n"+m.Style.Description.Render(m.Items[m.st.Cursor].Description()))
	return buf.String()
}

func (m Model) editView() string {
	item := m.Items[m.st.Cursor]

	var v string
	switch m.edittype {
	case TText:
		v = m.textinput.View()
	case TMultiline:
		v = m.textarea.View()
	case TRadio:
		v = m.radio.View()
	case TFileExisting:
		v = m.filemgr.View()
	default:
		return "INTERNAL ERROR"
	}
	return "--[" + item.Name() + "]------\n" + v + "\n" + m.Style.Description.Render(item.Description())
}
