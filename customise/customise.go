package customise

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"bbtea/display"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Items     []Item
	cursor    int
	width     int
	editing   bool
	edittype  Type
	finishing bool
	err       error
	fields
}

type fields struct {
	textarea  textarea.Model
	textinput textinput.Model
	radio     RadioButton
}

func NewModel(items []Item) Model {
	return Model{
		Items: items,
		fields: fields{
			textarea:  textarea.New(),
			textinput: textinput.New(),
			radio:     RadioButton{},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.editing {
		return m.procMsgEdit(msg)
	}
	return m.procMsgView(msg)

}
func (m Model) procMsgView(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
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
		case " ":
			if m.Items[m.cursor].Type() != TCheckbox {
				break
			}
			fallthrough
		case "enter", "f4":
			item := m.Items[m.cursor]

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
			case TCheckbox:
				if item.Value() == sTrue {
					item.Set(sFalse)
				} else {
					item.Set(sTrue)
				}
			}
		}
	}
	return m, nil
}

func (m Model) procMsgEdit(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// we only process enter for non-multiline modes
			if m.edittype == TMultiline {
				break
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
			}
			m.Items[m.cursor].Set(val)
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
		var val string
		switch item.Type() {
		case TMultiline, TText:
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
		fmt.Fprintf(tw, "%s%s\t%v\n", cursor, item.Name(), val)
	}
	tw.Flush()

	// description
	fmt.Fprint(&buf, "\n"+m.Items[m.cursor].Description())
	return buf.String()
}

func (m Model) editView() string {
	field := m.Items[m.cursor]

	var v string
	switch m.edittype {
	case TText:
		v = m.textinput.View()
	case TMultiline:
		v = m.textarea.View()
	case TRadio:
		v = m.radio.View()
	default:
		return "INTERNAL ERROR"
	}
	return "--[" + field.Name() + "]------\n" + v + "\n\n" + field.Description()
}
