package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/models"
	"time"
)

type AddClients struct {}

type ChangeViewMsg struct {
	newView View
	params  []interface{}
}

func ChangeViewCmd(view View, params... interface{}) func() tea.Msg {
	return func() tea.Msg {
		return ChangeViewMsg{view, params}
	}
}

type ErrorMsg struct {
	err error
}

func ErrorCmd(err error) func() tea.Msg {
	return func() tea.Msg {
		return ErrorMsg{err}
	}
}

func login(username, password string) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		rtp, err := client.Login(username, password)
		if err != nil {
			return ErrorMsg{err}
		}
		client.tp = *rtp

		return ChangeViewMsg{ViewUser, nil}
	}
}

func register(username, email, password string) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		err = client.Register(username, password, email)
		if err != nil {
			return ErrorMsg{err}
		}
		return ChangeViewMsg{ViewHome, nil}
	}
}

type UserMsg struct {
	user models.User
	aurumAdmin bool
}

func getUser() tea.Msg {
	client, err := clientManager.GetActiveClient()
	if err != nil {
		return ErrorMsg{err}
	}

	user, err := client.GetUserInfo(&client.tp)
	if err != nil {
		return ErrorMsg{err}
	}

	admin, err := client.GetAccess(aurum.AurumName, user.Username)
	if err != nil {
		return ErrorMsg{err}
	}

	return UserMsg{
		*user,
		admin.AllowedAccess && admin.Role == models.RoleAdmin,
	}
}


type UsersMsg struct {
	users []models.User
}

func getUsers() tea.Msg {
	client, err := clientManager.GetActiveClient()
	if err != nil {
		return ErrorMsg{err}
	}

	users, err := client.GetUsers(&client.tp)
	if err != nil {
		return ErrorMsg{err}
	}

	return UsersMsg{
		users,
	}
}


type GroupsMsg struct {
	groups []models.Group
}

func getGroups() tea.Msg {
	client, err := clientManager.GetActiveClient()
	if err != nil {
		return ErrorMsg{err}
	}

	groups, err := client.GetGroups()
	if err != nil {
		return ErrorMsg{err}
	}

	return GroupsMsg {
		groups,
	}
}


type UserGroupsMsg struct {
	groups []models.GroupWithRole
	index int
}

func getGroupsForUser(user string, index int) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		groups, err := client.GetGroupsForUser(&client.tp, user)
		if err != nil {
			return ErrorMsg{err}
		}

		return UserGroupsMsg{
			groups,
			index,
		}
	}
}

type UpdateUserMsg struct {
	user *models.User
}

func changePassword(password string) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		_, err = client.UpdateUser(&client.tp, &models.User{
			Password: password,
		})
		if err != nil {
			return ErrorMsg{err}
		}

		return Compound(ChangeViewCmd(ViewUser)(), getUser())
	}
}

func changeEmail(email string) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		_, err = client.UpdateUser(&client.tp, &models.User{
			Email: email,
		})
		if err != nil {
			return ErrorMsg{err}
		}

		return Compound(ChangeViewCmd(ViewUser)(), getUser())
	}
}

func newGroup(name string) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		err = client.AddGroup(&client.tp, &models.Group{
			Name: name,
			AllowRegistration: true,
		})
		if err != nil {
			return ErrorMsg{err}
		}

		return getGroups()
	}
}

func deleteGroup(name string) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		err = client.RemoveGroup(&client.tp, name)
		if err != nil {
			return ErrorMsg{err}
		}

		time.Sleep(2 * time.Second)
		return getGroups()
	}
}

func addUserToGroup(username, groupname string) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		err = client.AddUserToGroup(&client.tp, username, groupname)
		if err != nil {
			return ErrorMsg{err}
		}

		return getGroupsForUser(username, -1)
	}
}


func removeUserFromGroup(username, groupname string) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		err = client.RemoveUserFromGroup(&client.tp, username, groupname)
		if err != nil {
			return ErrorMsg{err}
		}

		return getGroupsForUser(username, -1)
	}
}


func setAccess(username, groupname string, role models.Role) func() tea.Msg {
	return func() tea.Msg {
		client, err := clientManager.GetActiveClient()
		if err != nil {
			return ErrorMsg{err}
		}

		err = client.SetAccess(&client.tp, models.AccessStatus{
			GroupName:     groupname,
			Username:      username,
			Role:          role,
			AllowedAccess: true,
		})
		if err != nil {
			return ErrorMsg{err}
		}

		return getGroupsForUser(username, -1)
	}
}


type CompoundMsg struct {
	msgs []tea.Msg
}

func Compound(msgs... tea.Msg) tea.Msg {
	return CompoundMsg{
		msgs: msgs,
	}
}