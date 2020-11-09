package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/pkg/models"
	te "github.com/muesli/termenv"
)

type UserModel struct {
	user *models.User
}

func InitialUserScreenModel() UserModel {
	return UserModel{&models.User{}}
}

func (m UserModel) View() string {
	username := te.String(m.user.Username).Foreground(color("#00f")).String()

	s := fmt.Sprintf(" Welcome %s\n\n", username)

	return s
}

func (m UserModel) Update(msg tea.Msg) (UserModel, tea.Cmd) {
	switch msg := msg.(type) {
	case getMeMsg:
		m.user = msg.user
	}

	return m, nil
}
