package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	te "github.com/muesli/termenv"
	"strings"
)

const focusedTextColor = "205"

var (
	focusedPrompt       = te.String("> ").Foreground(color(focusedTextColor)).String()
	blurredPrompt       = "> "
	focusedSubmitButton = "[ " + te.String("Submit").Foreground(color(focusedTextColor)).String() + " ]"
	blurredSubmitButton = "[ " + te.String("Submit").Foreground(color("240")).String() + " ]"
)

type LoginScreenModel struct {
	index         int
	usernameInput textinput.Model
	passwordInput textinput.Model
	emailInput    textinput.Model
	submitButton  string
	// login determines if to show login or register screen
	login bool

	err error
}

func initialLoginScreenModel() LoginScreenModel {
	name := textinput.NewModel()
	name.Placeholder = "username"
	name.Focus()
	name.Prompt = focusedPrompt
	name.TextColor = focusedTextColor

	email := textinput.NewModel()
	email.Placeholder = "example@email.com"
	email.Prompt = blurredPrompt

	password := textinput.NewModel()
	password.Placeholder = "password"
	password.Prompt = blurredPrompt
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	return LoginScreenModel{0, name, password, email, blurredSubmitButton, true, nil}
}

func LoginScreenView(m model) string {
	var header string
	if m.login.login {
		header = "Login"
	} else {
		header = "Register"
	}

	s := fmt.Sprintf(" %s %s\n", aurumText, header)

	if m.login.err != nil {
		s += te.String("Error: ").Foreground(color("#f00")).String() + strings.TrimSpace(m.login.err.Error()) + "\n"
	}

	s += "\n"

	var inputs []string
	inputs = append(inputs, m.login.usernameInput.View())
	if !m.login.login {
		inputs = append(inputs, m.login.emailInput.View())
	}
	inputs = append(inputs, m.login.passwordInput.View())

	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n\n" + m.login.submitButton + "\n"
	return s
}

func LoginScreenMsgHandler(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case loginErrMsg:
		m.login.err = msg.err
	case tea.KeyMsg:
		switch key := msg.Type; key {
		case tea.KeyShiftTab, tea.KeyTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			var inputs []textinput.Model
			inputs = append(inputs, m.login.usernameInput)
			if !m.login.login {
				inputs = append(inputs, m.login.emailInput)
			}
			inputs = append(inputs, m.login.passwordInput)

			// if submit login
			if key == tea.KeyEnter && m.login.index == len(inputs) {
				if m.login.login {
					return m, login(m.au, m.login.usernameInput.Value(), m.login.passwordInput.Value())
				} else {
					return m, register(m.au, m.login.usernameInput.Value(), m.login.emailInput.Value(), m.login.passwordInput.Value())
				}
			}

			if key == tea.KeyUp || key == tea.KeyShiftTab {
				m.login.index--
			} else {
				m.login.index++
			}

			if m.login.index > len(inputs) {
				m.login.index = 0
			} else if m.login.index < 0 {
				m.login.index = len(inputs)
			}

			for i := 0; i <= len(inputs)-1; i++ {
				if i == m.login.index {
					// Set focused state
					inputs[i].Focus()
					inputs[i].Prompt = focusedPrompt
					inputs[i].TextColor = focusedTextColor
					continue
				}
				// Remove focused state
				inputs[i].Blur()
				inputs[i].Prompt = blurredPrompt
				inputs[i].TextColor = ""
			}

			m.login.usernameInput = inputs[0]
			if m.login.login {
				m.login.passwordInput = inputs[1]
			} else {
				m.login.emailInput = inputs[1]
				m.login.passwordInput = inputs[2]
			}

			if m.login.index == len(inputs) {
				m.login.submitButton = focusedSubmitButton
			} else {
				m.login.submitButton = blurredSubmitButton
			}

			return m, nil
		}
	}

	m.login, cmd = updateInputs(msg, m.login)
	return m, cmd
}

// Pass messages and models through to text input components. Only text inputs
// with Focus() set will respond, so it's safe to simply update all of them
// here without any further logic.
func updateInputs(msg tea.Msg, m LoginScreenModel) (LoginScreenModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.usernameInput, cmd = m.usernameInput.Update(msg)
	cmds = append(cmds, cmd)

	m.emailInput, cmd = m.emailInput.Update(msg)
	cmds = append(cmds, cmd)

	m.passwordInput, cmd = m.passwordInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
