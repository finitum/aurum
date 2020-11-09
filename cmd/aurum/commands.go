package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/clients/goclient"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
)

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

type connectMsg struct {
	au *goclient.Aurum
}

type loginMsg struct {
	tp *jwt.TokenPair
}

type getMeMsg struct {
	user *models.User
}

type loginErrMsg struct {
	err error
}

type registerMsg struct{}

func connect() tea.Msg {
	au, err := goclient.Connect(host)
	if err != nil {
		return errMsg{err}
	}

	return connectMsg{au}
}

func login(au *goclient.Aurum, username, password string) tea.Cmd {
	return func() tea.Msg {
		tp, err := au.Login(username, password)
		if err != nil {
			return loginErrMsg{err}
		}

		return loginMsg{tp}
	}
}

func register(au *goclient.Aurum, username, email, password string) tea.Cmd {
	return func() tea.Msg {
		err := au.Register(username, password, email)
		if err != nil {
			return loginErrMsg{err}
		}

		return registerMsg{}
	}
}

func getme(au *goclient.Aurum, tp *jwt.TokenPair) tea.Cmd {
	return func() tea.Msg {
		user, err := au.GetUserInfo(tp)
		if err != nil {
			return errMsg{err}
		}

		return getMeMsg{user}
	}
}
