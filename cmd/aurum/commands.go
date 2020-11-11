package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/clients/go"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/oldmodels"
)

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

type connectMsg struct {
	au *aurum.Aurum
}

type loginMsg struct {
	tp *jwt.TokenPair
}

type loginErrMsg struct {
	err error
}

type registerMsg struct{}

type getMeMsg struct {
	user *oldmodels.User
}

func connect() tea.Msg {
	au, err := aurum.Connect(*host)
	if err != nil {
		return errMsg{err}
	}

	return connectMsg{au}
}

func getme(au *aurum.Aurum, tp *jwt.TokenPair) tea.Cmd {
	return func() tea.Msg {
		user, err := au.GetUserInfo(tp)
		if err != nil {
			return errMsg{err}
		}

		return getMeMsg{user}
	}
}

func login(au *aurum.Aurum, username, password string) tea.Cmd {
	return func() tea.Msg {
		tp, err := au.Login(username, password)
		if err != nil {
			return loginErrMsg{err}
		}

		return loginMsg{tp}
	}
}

func register(au *aurum.Aurum, username, email, password string) tea.Cmd {
	return func() tea.Msg {
		err := au.Register(username, password, email)
		if err != nil {
			return loginErrMsg{err}
		}

		return registerMsg{}
	}
}
