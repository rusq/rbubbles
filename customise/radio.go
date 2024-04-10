package customise

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type RadioButton struct {
	choices  []string
	selected int
}

func (r RadioButton) Init() tea.Cmd {
	return nil
}

func (r *RadioButton) SetValues(v []string, selected string) {
	r.choices = v
	r.selected = 0
	for i := range v {
		if strings.EqualFold(selected, v[i]) {
			r.selected = i
			break
		}
	}
}

func (r RadioButton) Update(msg tea.Msg) (RadioButton, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if 0 < r.selected {
				r.selected--
			}
		case "down":
			if r.selected < len(r.choices)-1 {
				r.selected++
			}
		}
	}
	return r, nil
}

func (r RadioButton) View() string {
	var buf strings.Builder
	for i, val := range r.choices {
		cur := "( ) "
		if i == r.selected {
			cur = "(*) "
		}
		fmt.Fprintf(&buf, "%s%s\n", cur, val)
	}
	return buf.String()
}

func (r RadioButton) Value() string {
	return r.choices[r.selected]
}
