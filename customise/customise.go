package customise

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"bbtea/display"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Items     []Item
	cursor    int
	width     int
	editing   bool
	editbox   textarea.Model
	finishing bool
	err       error
	errField  int
}

func NewModel(items []Item) Model {
	return Model{
		Items:   items,
		editbox: textarea.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "esc":
				m.editing = false
				m.editbox.Blur()
				m.Items[m.cursor].Set(m.editbox.Value())
			}
		} else {
			switch msg.String() {
			case "q":
				m.finishing = true
				return m, tea.Quit
			case "j", "down":
				if m.cursor < len(m.Items)-1 {
					m.cursor++
				}
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}

			case "enter":
				if !m.editing {
					m.editing = !m.editing
					m.editbox.SetValue(m.Items[m.cursor].Value())
					m.editbox.Focus()
				}
			}
		}
	}
	var cmds []tea.Cmd
	if m.editing {
		var cmd tea.Cmd
		m.editbox, cmd = m.editbox.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.finishing {
		return ""
	}
	if m.err != nil {
		return m.err.Error()
	}
	if m.editing {
		return m.editbox.View()
	}
	if len(m.Items) == 0 {
		return "No items to show."
	}

	var buf strings.Builder
	tw := tabwriter.NewWriter(&buf, 0, 4, 4, ' ', 0)

	for i, item := range m.Items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		value := item.Value()
		if len(value) == 0 {
			value = "<empty>"
		}
		fmt.Fprintf(tw, "%s%s\t%v\n", cursor, item.Name(), display.Trunc(value, m.width))
		// fmt.Fprintf(&buf, " %s\n",  item.Description())
	}
	tw.Flush()
	return buf.String()
}
