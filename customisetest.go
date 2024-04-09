package main

import (
	"bbtea/customise"

	tea "github.com/charmbracelet/bubbletea"
)

func customiseTest() {
	var testVar string = "Hello, World!"
	var testInt int = 42

	c := customise.NewModel([]customise.Item{
		customise.StringVar(&testVar, "TestVar", "This is a test variable", "Test"),
		customise.IntVar(&testInt, "TestInt", "This is a test integer", "Test"),
	})
	p := tea.NewProgram(custmodel{c})
	_, err := p.Run()
	if err != nil {
		panic(err)
	}
}

type custmodel struct {
	m customise.Model
}

func (f custmodel) Init() tea.Cmd {
	return f.m.Init()
}

func (f custmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	f.m, cmd = f.m.Update(msg)
	return f, cmd
}

func (f custmodel) View() string {
	return f.m.View()
}
