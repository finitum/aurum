package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/clients/goclient"
	"github.com/finitum/aurum/pkg/jwt"
	te "github.com/muesli/termenv"
	"os"
)

const host = "http://localhost:8042"

var (
	color     = te.ColorProfile().Color
	aurumText = te.String("Aurum").Foreground(color("#ffd700")).Bold().String()
)

func main() {
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
	au     *goclient.Aurum
	tp     *jwt.TokenPair
	screen Screen
	err    error

	info string

	main  MainScreenModel
	login LoginScreenModel
	user  UserScreenModel
}

func initialModel() model {
	return model{
		au:     nil,
		screen: MainScreen,
		err:    nil,
		main:   initialMainScreenModel(),
		login:  initialLoginScreenModel(),
		user:   initialUserScreenModel(),
	}
}

func (m model) Init() tea.Cmd {
	return connect
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case getMeMsg:
		m.user.user = msg.user
	case loginMsg:
		m.screen = UserScreen
		m.login = initialLoginScreenModel()
		m.tp = msg.tp
		return m, getme(m.au, msg.tp)
	case registerMsg:
		m.info = te.String("Registered successfully!").Foreground(color("#0f0")).String()
		m.screen = MainScreen
		return m, nil
	case errMsg:
		m.err = msg
		return m, tea.Quit
	case connectMsg:
		m.au = msg.au
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			m.screen = MainScreen
			return m, nil
		}
	}

	switch m.screen {
	case MainScreen:
		return MainScreenMsgHandler(m, msg)
	case LoginScreen, RegisterScreen:
		return LoginScreenMsgHandler(m, msg)
	case UserScreen:
		return UserScreenMsgHandler(m, msg)
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nAn error occured: %v\n\n", m.err)
	}

	if m.au == nil {
		return fmt.Sprint("Connecting...\n")
	}

	switch m.screen {
	default:
		fallthrough
	case MainScreen:
		return MainScreenView(m)
	case LoginScreen, RegisterScreen:
		return LoginScreenView(m)
	case UserScreen:
		return UserScreenView(m)
	}
}
