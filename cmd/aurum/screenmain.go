package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MainScreenModel struct {
	choices []string
	cursor  int
}

func initialMainScreenModel() MainScreenModel {
	return MainScreenModel{
		cursor: 0,
		choices: []string{
			"Login",
			"Register",
			"Exit",
		},
	}
}

func MainScreenView(m model) string {
	s := fmt.Sprintf(" %s \n", aurumText)

	if m.info != "" {
		s += m.info + "\n"
	}

	s += "\n"

	for i, choice := range m.main.choices {
		cursor := " "
		if i == m.main.cursor {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s
}

func MainScreenMsgHandler(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.info = ""
		switch msg.Type {
		case tea.KeyUp:
			m.main.cursor--
			if m.main.cursor < 0 {
				m.main.cursor = len(m.main.choices) - 1
			}
		case tea.KeyDown:
			m.main.cursor++
			if m.main.cursor > len(m.main.choices)-1 {
				m.main.cursor = 0
			}
		case tea.KeyEnter:
			switch m.main.choices[m.main.cursor] {
			case "Login":
				m.screen = LoginScreen
				m.login.login = true
				return m, textinput.Blink
			case "Register":
				m.login.login = false
				m.screen = RegisterScreen
				return m, textinput.Blink
			case "Exit":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}
