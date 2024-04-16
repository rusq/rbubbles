package main

import (
	"fmt"
	"os"

	"github.com/rusq/rbubbles/filemgr"

	tea "github.com/charmbracelet/bubbletea"
)

func filebrowser() {
	fm := filemgr.New(os.DirFS("."), ".", 10, "*")
	fm.Focus()
	// fm.ShowHelp = true
	fm.Debug = os.Getenv("DEBUG") != ""
	p := tea.NewProgram(fmmodel{fm, false})
	r, err := p.Run()
	if err != nil {
		panic(err)
	}
	m := r.(fmmodel)
	fmt.Println(m.m.Selected)
}

type fmmodel struct {
	m         filemgr.Model
	finishing bool
}

func (f fmmodel) Init() tea.Cmd {
	return f.m.Init()
}

func (f fmmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			f.finishing = true
			return f, tea.Quit
		}
	}
	f.m, cmd = f.m.Update(msg)
	return f, cmd
}

func (f fmmodel) View() string {
	if f.finishing {
		return ""
	}
	return f.m.View()
}
