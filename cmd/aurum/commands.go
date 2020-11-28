package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/clients/go"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
)

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

type connectMsg struct {
	au aurum.Client
}

type loginMsg struct {
	tp *jwt.TokenPair
}

type loginErrMsg struct {
	err error
}

type registerMsg struct{}

type getMeMsg struct {
	user *models.User
}

func connect() tea.Msg {
	au, err := aurum.NewRemoteClient(*host)
	if err != nil {
		return errMsg{err}
	}

	return connectMsg{au}
}

func getme(au aurum.Client, tp *jwt.TokenPair) tea.Cmd {
	return func() tea.Msg {
		user, err := au.GetUserInfo(tp)
		if err != nil {
			return errMsg{err}
		}

		return getMeMsg{user}
	}
}

func login(au aurum.Client, username, password string) tea.Cmd {
	return func() tea.Msg {
		tp, err := au.Login(username, password)
		if err != nil {
			return loginErrMsg{err}
		}

		return loginMsg{tp}
	}
}

func register(au aurum.Client, username, email, password string) tea.Cmd {
	return func() tea.Msg {
		err := au.Register(username, password, email)
		if err != nil {
			return loginErrMsg{err}
		}

		return registerMsg{}
	}
}
