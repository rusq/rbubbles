package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func wiztest() {
	m := NewWizModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		log.Fatalf("Error starting program: %v", err)
	}
	log.Print(m.form.GetString("choice"))
}

type wizModel struct {
	form *huh.Form
	help []string
}

func NewWizModel() *wizModel {
	return &wizModel{
		help: []string{"archive conversations", "list users or channels", "dump conversations"},
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("choice").
					Options(huh.NewOptions("archive", "list", "dump")...).
					Title("Select an action"),
			),
		),
	}
}

func (w *wizModel) Init() tea.Cmd {
	return w.form.Init()
}

func (w *wizModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return w, tea.Quit
		}
	}

	var cmds []tea.Cmd
	form, cmd := w.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		w.form = f
		cmds = append(cmds, cmd)
	}

	if w.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return w, tea.Batch(cmds...)
}

func (w *wizModel) View() string {
	return w.form.View()
}
