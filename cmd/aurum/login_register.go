package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/clients/go"
	te "github.com/muesli/termenv"
)

const focusedTextColor = "205"

var (
	focusedPrompt       = te.String("> ").Foreground(color(focusedTextColor)).String()
	blurredPrompt       = "> "
	focusedSubmitButton = "[ " + te.String("Submit").Foreground(color(focusedTextColor)).String() + " ]"
	blurredSubmitButton = "[ " + te.String("Submit").Foreground(color("240")).String() + " ]"
)

type LoginRegisterModel struct {
	index         int
	usernameInput textinput.Model
	passwordInput textinput.Model
	emailInput    textinput.Model
	submitButton  string
	// login determines if to show login or register screen
	login bool

	err error
}

func InitialLoginScreenModel() LoginRegisterModel {
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

	return LoginRegisterModel{0, name, password, email, blurredSubmitButton, true, nil}
}

func (m LoginRegisterModel) View() string {
	var header string
	if m.login {
		header = "Login"
	} else {
		header = "Register"
	}

	s := fmt.Sprintf(" %s %s\n", aurumText, header)

	if m.err != nil {
		s += te.String("Error: ").Foreground(color("#f00")).String() + strings.TrimSpace(m.err.Error()) + "\n"
	}

	s += "\n"

	var inputs []string
	inputs = append(inputs, m.usernameInput.View())
	if !m.login {
		inputs = append(inputs, m.emailInput.View())
	}
	inputs = append(inputs, m.passwordInput.View())

	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n\n" + m.submitButton + "\n"
	return s
}

func (m LoginRegisterModel) Update(au aurum.Client, msg tea.Msg) (LoginRegisterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case loginErrMsg:
		m.err = msg.err
	case tea.KeyMsg:
		switch key := msg.Type; key {
		case tea.KeyShiftTab, tea.KeyTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			var inputs []textinput.Model
			inputs = append(inputs, m.usernameInput)
			if !m.login {
				inputs = append(inputs, m.emailInput)
			}
			inputs = append(inputs, m.passwordInput)

			// if submit login
			if key == tea.KeyEnter && m.index == len(inputs) {
				if m.login {
					return m, login(au, m.usernameInput.Value(), m.passwordInput.Value())
				} else {
					return m, register(au, m.usernameInput.Value(), m.emailInput.Value(), m.passwordInput.Value())
				}
			}

			if key == tea.KeyUp || key == tea.KeyShiftTab {
				m.index--
			} else {
				m.index++
			}

			if m.index > len(inputs) {
				m.index = 0
			} else if m.index < 0 {
				m.index = len(inputs)
			}

			for i := 0; i <= len(inputs)-1; i++ {
				if i == m.index {
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

			m.usernameInput = inputs[0]
			if m.login {
				m.passwordInput = inputs[1]
			} else {
				m.emailInput = inputs[1]
				m.passwordInput = inputs[2]
			}

			if m.index == len(inputs) {
				m.submitButton = focusedSubmitButton
			} else {
				m.submitButton = blurredSubmitButton
			}

			return m, nil
		}
	}

	return updateInputs(msg, m)
}

// Pass messages and models through to text input components. Only text inputs
// with Focus() set will respond, so it's safe to simply update all of them
// here without any further logic.
func updateInputs(msg tea.Msg, m LoginRegisterModel) (LoginRegisterModel, tea.Cmd) {
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
