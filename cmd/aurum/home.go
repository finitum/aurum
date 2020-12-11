package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

const exitChoice = "Exit"

var HomeChoices = []View{
	ViewLogin,
	ViewRegister,
	ViewChangeServer,
	View(exitChoice),
}


type HomeModel struct {
	cursor int
}

func NewHomeModel() HomeModel {
	return HomeModel {
		0,
	}
}

func (h HomeModel) View() string {
	s := "\n"
	
	for i, choice := range HomeChoices {
		cursor := " "
		if i == h.cursor {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	
	return s
}

func (h HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			h.cursor--
			if h.cursor < 0 {
				h.cursor = len(HomeChoices) - 1
			}
		case tea.KeyDown:
			h.cursor++
			if h.cursor > len(HomeChoices) - 1 {
				h.cursor = 0
			}
		case tea.KeyEnter:
			switch HomeChoices[h.cursor] {
			case exitChoice:
				return h, tea.Quit
			default:
				return h, ChangeViewCmd(HomeChoices[h.cursor])
			}
		}
	}

	return h, nil
}


