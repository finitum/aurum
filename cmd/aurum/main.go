package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/clients/go"
	"github.com/finitum/aurum/pkg/jwt"
	te "github.com/muesli/termenv"
	"os"
)

var host = flag.String("host", "http://localhost:8042", "Aurum host to connect to")

var (
	color     = te.ColorProfile().Color
	aurumText = te.String("Aurum").Foreground(color("#ffd700")).Bold().String()
)

func main() {
	flag.Parse()

	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type Screen int

const (
	MainScreen Screen = iota
	LoginScreen
	RegisterScreen
	UserScreen
)

type model struct {
	au     *aurum.Aurum
	tp     *jwt.TokenPair
	screen Screen
	err    error

	info string

	main  MainScreenModel
	login LoginRegisterModel
	user  UserModel
}

func initialModel() model {
	return model{
		au:     nil,
		screen: MainScreen,
		err:    nil,
		main:   InitialMainScreenModel(),
		login:  InitialLoginScreenModel(),
		user:   InitialUserScreenModel(),
	}
}

func (m model) Init() tea.Cmd {
	return connect
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case screenLoginMsg:
		m.screen = LoginScreen
		m.login = InitialLoginScreenModel()
		m.login.login = true
		m.screen = RegisterScreen
	case screenRegisterMsg:
		m.login = InitialLoginScreenModel()
		m.login.login = false
		m.screen = RegisterScreen
	case loginMsg:
		m.screen = UserScreen
		m.tp = msg.tp
		cmds = append(cmds, getme(m.au, msg.tp))
	case registerMsg:
		m.info = te.String("Registered successfully!").Foreground(color("#0f0")).String()
		m.screen = MainScreen
	case errMsg:
		m.err = msg
		return m, tea.Quit
	case connectMsg:
		m.au = msg.au
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			m.screen = MainScreen
		}
	}

	var cmd tea.Cmd
	switch m.screen {
	case MainScreen:
		m.main, cmd = m.main.Update(msg)
	case LoginScreen, RegisterScreen:
		m.login, cmd = m.login.Update(m.au, msg)
	case UserScreen:
		m.user, cmd = m.user.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nAn error occured: %v\n\n", m.err)
	}

	if m.au == nil {
		return fmt.Sprint("Connecting...\n")
	}

	var s string
	switch m.screen {
	case MainScreen:
		s += m.main.View()
	case LoginScreen, RegisterScreen:
		s += m.login.View()
	case UserScreen:
		s += m.user.View()
	}

	return s
}
