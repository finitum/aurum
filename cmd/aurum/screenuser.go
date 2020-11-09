package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/pkg/models"
	te "github.com/muesli/termenv"
)

type UserScreenModel struct {
	user *models.User
}

func initialUserScreenModel() UserScreenModel {
	return UserScreenModel{&models.User{}}
}

func UserScreenView(m model) string {
	username := te.String(m.user.user.Username).Foreground(color("#00f")).String()

	s := fmt.Sprintf(" Welcome %s\n\n", username)

	return s
}

func UserScreenMsgHandler(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
