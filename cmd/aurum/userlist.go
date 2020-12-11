package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/pkg/models"
	te "github.com/muesli/termenv"
	"strings"
)

type User struct {
	models.User
	groups []models.GroupWithRole
}

type UserListModel struct {
	cursor     int
	users      []User
	opened     int
	triedUsers bool
}

var header = te.String(" Username         Email address   %s      ViewGroupList\n").Bold()
var userRowFmtEven = te.String(" %%-15s  %%-%ds  %%-30s\n").Foreground(color("7"))
var userRowFmtOdd = te.String(" %%-15s  %%-%ds  %%-30s\n").Foreground(color("8"))
var userRowFmtSelected = te.String(" %%-15s  %%-%ds  %%-30s\n").Foreground(color(focusedTextColor))

const explanation = `
 [Enter] to edit user groups, 
 [Esc] to go back to the main menu
 [Up],[Down] to scroll

 (a): Admin in this group 
`

func NewUserListModel() UserListModel {
	return UserListModel{
		cursor:     0,
		users:      []User{},
		triedUsers: false,
	}
}

func (u UserListModel) Update(msg tea.Msg) (UserListModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			u.cursor--
			if u.cursor < 0 {
				u.cursor = len(u.users) - 1
			}
		case tea.KeyDown:
			u.cursor++
			if u.cursor > len(u.users)-1 {
				u.cursor = 0
			}
		case tea.KeyEnter:
			cmds = append(cmds, ChangeViewCmd(ViewEditUserGroups, &u.users[u.cursor].User))
		case tea.KeyEsc:
			cmds = append(cmds, ChangeViewCmd(ViewUser))
		}
	case UsersMsg:
		u.users = make([]User, len(msg.users))
		for index, i := range msg.users {
			u.users[index] = User{
				User:   i,
				groups: nil,
			}

			cmds = append(cmds, getGroupsForUser(i.Username, index))
		}
	case UserGroupsMsg:
		u.users[msg.index].groups = msg.groups
	}

	if !u.triedUsers {
		cmds = append(cmds, getUsers)
		u.triedUsers = true
	}

	return u, tea.Batch(cmds...)
}

func (u UserListModel) View(width int) string {
	s := "\n"

	s += explanation

	s += "\n"

	maxEmailLength := 20
	if width > 80 {
		maxEmailLength = 40
	}

	s += fmt.Sprintf(header.String(), strings.Repeat(" ", maxEmailLength-20)) + "\n"
	s += strings.Repeat("—", width) + "\n"
	for i, user := range u.users {
		email := user.Email
		if email == "" {
			email = "not set"
		}

		if len(email) > maxEmailLength {
			email = email[:maxEmailLength-2] + ".."
		}

		format := userRowFmtOdd.String()
		if i%2 == 0 {
			format = userRowFmtEven.String()
		}

		if i == u.cursor {
			format = userRowFmtSelected.String()
		}

		var groups string
		for index, i := range user.groups {
			groups += i.Name
			if i.Role == models.RoleAdmin {
				groups += "(a)"
			}

			if index != len(user.groups)-1 {
				groups += ","
			}
		}

		s += fmt.Sprintf(fmt.Sprintf(format, maxEmailLength), user.Username, email, groups)
	}

	return s
}

func (u UserListModel) Init(params []interface{}) (UserListModel, tea.Cmd)  {
	u.triedUsers = false
	return u, nil
}