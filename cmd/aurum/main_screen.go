package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MainScreenModel struct {
	choices []string
	cursor  int
}

func InitialMainScreenModel() MainScreenModel {
	return MainScreenModel{
		cursor: 0,
		choices: []string{
			"Login",
			"Register",
			"Exit",
		},
	}
}

func (m MainScreenModel) View() string {
	s := fmt.Sprintf(" %s \n", aurumText)
	s += "\n"

	for i, choice := range m.choices {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s
}

type screenLoginMsg struct{}

func toLoginScreen() tea.Msg { return screenLoginMsg{} }

type screenRegisterMsg struct{}

func toRegisterScreen() tea.Msg { return screenRegisterMsg{} }

func (m MainScreenModel) Update(msg tea.Msg) (MainScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		case tea.KeyDown:
			m.cursor++
			if m.cursor > len(m.choices)-1 {
				m.cursor = 0
			}
		case tea.KeyEnter:
			switch m.choices[m.cursor] {
			case "Login":
				return m, toLoginScreen
			case "Register":
				return m, toRegisterScreen
			case "Exit":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}
