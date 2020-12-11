package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	te "github.com/muesli/termenv"
	"go.deanishe.net/env"
	"strings"
)

var (
	color     = te.ColorProfile().Color
	aurumText = te.String("Aurum").Foreground(color("#ffd700")).Bold().String()
)


type View string
const (
	ViewHome           View = "Home"
	ViewLogin          View = "Login"
	ViewRegister       View = "Register"
	ViewUser           View = "User"
	ViewUserList       View = "User List"
	ViewGroupList      View = "Group list"
	ViewChangePassword View = "Change Password"
	ViewChangeEmail    View = "Change Email"
	ViewEditUserGroups View = "Edit User Groups"
	ViewChangeServer   View = "Change Server"
)

type Model struct {
	login          LoginModel
	home           HomeModel
	register       RegisterModel
	user           UserModel
	userlist       UserListModel
	changePassword ChangePasswordModel
	changeEmail    ChangeEmailModel
	groupList      GroupListModel
	editUserGroups EditUserGroupModel
	changeServer   ChangeServerModel

	currentView View

	width int
	err error
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return AddClients{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case ErrorMsg:
		m.err = msg.err
	case ChangeViewMsg:
		prevview := m.currentView
		m.currentView = msg.newView
		switch m.currentView {
		case ViewHome:
		case ViewLogin:
		case ViewRegister:
		case ViewUser:
		case ViewUserList:
			m.userlist , cmd= m.userlist.Init(msg.params)
			cmds = append(cmds, cmd)
		case ViewChangePassword:
		case ViewChangeEmail:
		case ViewGroupList:
		case ViewEditUserGroups:
			m.editUserGroups, cmd = m.editUserGroups.Init(msg.params)
			cmds = append(cmds, cmd)
		case ViewChangeServer:
			m.changeServer, cmd = m.changeServer.Init([]interface{}{prevview})
			cmds = append(cmds, cmd)
		}

	case CompoundMsg:
		var newm tea.Model
		newm = m
		for _, s := range msg.msgs {
			newm, cmd = newm.Update(s)
			cmds = append(cmds, cmd)
		}

		return newm, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.err != nil {
				m.err = nil
				return m, nil
			}
		case tea.KeyCtrlS:
			prevview := m.currentView
			m.currentView = ViewChangeServer
			m.changeServer, cmd = m.changeServer.Init([]interface{}{prevview})
		}
	case AddClients:
		hostvar := env.Get("AURUM_TUI_HOST")
		if hostvar != "" {
			for _, i := range strings.Split(hostvar, ",") {
				err := clientManager.AddClient(i)
				if err != nil {
					return m, ErrorCmd(fmt.Errorf("couldn't connect with client %s (%v)", i, err))
				}
			}
		}

		for _, i := range hostFlags {
			err := clientManager.AddClient(i)
			if err != nil {
				return m, ErrorCmd(fmt.Errorf("couldn't connect with client %s (%v)", i, err))
			}
		}
	}

	switch m.currentView {
	case ViewHome:
		m.home, cmd = m.home.Update(msg)
	case ViewLogin:
		m.login, cmd = m.login.Update(msg)
	case ViewRegister:
		m.register, cmd = m.register.Update(msg)
	case ViewUser:
		m.user, cmd = m.user.Update(msg)
	case ViewUserList:
		m.userlist, cmd = m.userlist.Update(msg)
	case ViewChangePassword:
		m.changePassword, cmd = m.changePassword.Update(msg)
	case ViewChangeEmail:
		m.changeEmail, cmd = m.changeEmail.Update(msg)
	case ViewGroupList:
		m.groupList, cmd = m.groupList.Update(msg)
	case ViewEditUserGroups:
		m.editUserGroups, cmd = m.editUserGroups.Update(msg)
	case ViewChangeServer:
		m.changeServer, cmd = m.changeServer.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}


func (m Model) View() string {
	s := ""

	client, err := clientManager.GetActiveClient()
	var host string
	if err != nil {
		host = "<no host connected>"
	} else {
		host = client.GetUrl()
	}

	s += " " + aurumText + " at " + host + "\n"
	s += " press [control s] to change Aurum server\n"
	if m.width != -1 {
		s += strings.Repeat("=", m.width) + "\n"
	}

	var screen string
	switch m.currentView {
	case ViewLogin:
		screen = m.login.View()
	case ViewRegister:
		screen = m.register.View()
	case ViewHome:
		screen = m.home.View()
	case ViewUser:
		screen = m.user.View()
	case ViewUserList:
		screen = m.userlist.View(m.width)
	case ViewChangePassword:
		screen = m.changePassword.View()
	case ViewChangeEmail:
		screen = m.changeEmail.View()
	case ViewGroupList:
		screen = m.groupList.View(m.width)
	case ViewEditUserGroups:
		screen = m.editUserGroups.View(m.width)
	case ViewChangeServer:
		screen = m.changeServer.View()
	}


	for _, i := range strings.Split(screen, "\n") {
		if  m.width != -1 && runewidth.StringWidth(i) >= m.width {
			s += runewidth.Truncate(i, m.width, "") + "\n"
		} else {
			s += i + "\n"
		}
	}

	// Show last error
	if m.err != nil {
		s += te.String("\nError: ").Foreground(color("#f00")).String() + strings.TrimSpace(m.err.Error()) + "\n"
	}

	return s
}

func NewModel() Model {
	return Model{
		currentView:    ViewHome,
		width: 			-1,
		login:          NewLoginModel(),
		register:       NewRegisterModel(),
		home:           NewHomeModel(),
		user:           NewUserModel(),
		changePassword: NewChangePasswordModel(),
		changeEmail:    NewChangeEmailModel(),
		groupList: 		NewGroupListModel(),
		editUserGroups: NewEditUserGroupModel(),
		changeServer: 	NewChangeServerModel(),
	}
}

var clientManager = NewClientManager()

type arrayFlags []string
func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var hostFlags arrayFlags

func main() {
	flag.Var(&hostFlags, "host", "A host to connect to (multiple values possible)")
	flag.Parse()

	p := tea.NewProgram(NewModel())
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		fmt.Printf("Error starting Aurum TUI %v", err)
	}
}
