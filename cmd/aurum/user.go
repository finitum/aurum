package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
)

const logoutChoice = View("Logout")

var UserChoices = []View{
	ViewChangePassword,
	ViewChangeEmail,
	logoutChoice,
}

var AdminChoices = append(UserChoices, ViewUserList, ViewGroupList)

type UserModel struct {
	cursor int
	user *models.User
	admin bool
	triedUser bool
}

func NewUserModel() UserModel {
	return UserModel {
		0,
		nil,
		false,
		false,
	}
}

func(u UserModel) getChoices() []View {
	if u.admin {
		return AdminChoices
	} else {
		return UserChoices
	}
}

func (u UserModel) View() string {
	s := "\n"

	if u.user == nil {
		s += "\nWelcome"
	} else {
		s += "\nWelcome " + u.user.Username
	}

	s += "\n"

	for i, choice := range u.getChoices() {
		cursor := " "
		if i == u.cursor {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s
}

func (u UserModel) Update(msg tea.Msg) (UserModel, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			// Logout
			client, err := clientManager.GetActiveClient()
			if err != nil {
				return u, ErrorCmd(err)
			}
			client.tp = jwt.TokenPair{}
			return NewUserModel(), ChangeViewCmd(ViewHome)
		case tea.KeyUp:
			u.cursor--
			if u.cursor < 0 {
				u.cursor = len(u.getChoices()) - 1
			}
		case tea.KeyDown:
			u.cursor++
			if u.cursor > len(u.getChoices()) - 1 {
				u.cursor = 0
			}
		case tea.KeyEnter:
			switch u.getChoices()[u.cursor] {
			case logoutChoice:
				// Logout
				client, err := clientManager.GetActiveClient()
				if err != nil {
					return u, ErrorCmd(err)
				}
				client.tp = jwt.TokenPair{}
				return NewUserModel(), ChangeViewCmd(ViewHome)
			default:
				return u, ChangeViewCmd(u.getChoices()[u.cursor])
			}
		}
	case UserMsg:
		u.user = &msg.user
		u.admin = msg.aurumAdmin
	}

	if !u.triedUser {
		cmds = append(cmds, getUser)
		u.triedUser = true
	}

	return u, tea.Batch(cmds...)
}


