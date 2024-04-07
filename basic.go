package main

import tea "github.com/charmbracelet/bubbletea"

type model struct {
	choices  []string
	cursor   int
	selected map[string]struct{}
}

func initalModel() model {
	return model{
		choices:  []string{"one", "two", "three", "four", "five"},
		selected: make(map[string]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "ctrl+n":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.choices[m.cursor]]
			if ok {
				delete(m.selected, m.choices[m.cursor])
			} else {
				m.selected[m.choices[m.cursor]] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "pick some items\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[choice]; ok {
			checked = "x"
		}

		s += cursor + " " + checked + " " + choice + "\n"
	}

	s += "\nPress q to quit.\n"
	return s
}
