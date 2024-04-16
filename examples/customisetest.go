package main

import (
	"github.com/rusq/rbubbles/customise"

	tea "github.com/charmbracelet/bubbletea"
)

func customiseTest() {
	var testVar string = "Hello, World!"
	var testInt int = 42
	var testMultiline string = "Hello world\nMultiline"
	var testBool bool = true
	var testRadio string = "foo"
	var testFilename = "check_url.go"

	c := customise.NewModel([]customise.Item{
		customise.StringVar(&testVar, "TestVar", "This is a test variable", "Test"),
		customise.IntVar(&testInt, "TestInt", "This is a test integer", "Test"),
		customise.MultilineVar(&testMultiline, "Multiline test", "This is multiline test string", "Test"),
		customise.BoolVar(&testBool, "Boolean test", "This is boolean(checkbox) test", "Test"),
		customise.RadioStringVar(&testRadio, "test choice", "This is test choice", "Test", []string{"foo", "bar"}),
		customise.FilenameVar(&testFilename, "Filename test", "This is filename test", "Test", true),
	})
	p := tea.NewProgram(custmodel{m: c})
	_, err := p.Run()
	if err != nil {
		panic(err)
	}
}

type custmodel struct {
	m         customise.Model
	finishing bool
}

func (m custmodel) Init() tea.Cmd {
	return m.m.Init()
}

func (m custmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.finishing = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.m, cmd = m.m.Update(msg)
	return m, cmd
}

func (m custmodel) View() string {
	if m.finishing {
		return ""
	}
	return m.m.View()
}
