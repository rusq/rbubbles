package main

import (
	"fmt"
	"os"

	"bbtea/filemgr"

	tea "github.com/charmbracelet/bubbletea"
)

func filebrowser() {
	fm := filemgr.NewModel(".", 20, "*")
	fm.Debug = os.Getenv("DEBUG") != ""
	p := tea.NewProgram(fmmodel{fm})
	r, err := p.Run()
	if err != nil {
		panic(err)
	}
	m := r.(fmmodel)
	fmt.Println(m.m.Selected)
}

type fmmodel struct {
	m filemgr.Model
}

func (f fmmodel) Init() tea.Cmd {
	return f.m.Init()
}

func (f fmmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	f.m, cmd = f.m.Update(msg)
	return f, cmd
}

func (f fmmodel) View() string {
	return f.m.View()
}
