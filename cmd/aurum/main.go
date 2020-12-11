package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	aurum "github.com/finitum/aurum/clients/go"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/mattn/go-runewidth"
	te "github.com/muesli/termenv"
	"go.deanishe.net/env"
	"log"
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

	currentView View

	width int
	err error
}

var hostdefault = "http://localhost:8042"
var host = flag.String("host", hostdefault, "Aurum host to connect to")

var client aurum.Client
var tp jwt.TokenPair

func (m Model) Init() tea.Cmd {
	return nil
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
		m.currentView = msg.newView
		if msg.newView == ViewUserList {
			m.userlist.triedUsers = false
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
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}


func (m Model) View() string {
	s := " " + aurumText + " at " + *host

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
	}
}

func main() {
	hostvar := env.Get("AURUM_TUI_HOST")
	if env.Get("AURUM_TUI_HOST") != "" {
		hostdefault = hostvar
	}

	flag.Parse()

	var err error
	client, err = aurum.NewRemoteClient(*host)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(NewModel())
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		fmt.Printf("Error starting Aurum TUI %v", err)
	}
}
