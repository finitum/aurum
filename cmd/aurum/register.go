package main

import (
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type RegisterModel struct {
	index         int
	usernameInput textinput.Model
	emailInput	  textinput.Model
	passwordInput textinput.Model
	passwordRepeatInput textinput.Model
	submitButton  string
}

func NewRegisterModel() RegisterModel {
	name := textinput.NewModel()
	name.Placeholder = "username"
	name.Focus()
	name.Prompt = focusedPrompt
	name.TextColor = focusedTextColor

	email := textinput.NewModel()
	email.Placeholder = "email"
	email.Prompt = focusedPrompt
	email.TextColor = focusedTextColor

	password := textinput.NewModel()
	password.Placeholder = "password"
	password.Prompt = blurredPrompt
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	passwordrepeat := textinput.NewModel()
	passwordrepeat.Placeholder = "repeat password"
	passwordrepeat.Prompt = blurredPrompt
	passwordrepeat.EchoMode = textinput.EchoPassword
	passwordrepeat.EchoCharacter = '•'

	return RegisterModel{
		index:         0,
		usernameInput: name,
		emailInput: email,
		passwordInput: password,
		passwordRepeatInput: passwordrepeat,
		submitButton:  blurredSubmitButton,
	}
}

func (l RegisterModel) Update(msg tea.Msg) (RegisterModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
		case tea.KeyMsg:
		switch key := msg.Type; key {
		case tea.KeyShiftTab, tea.KeyTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			inputs := []textinput.Model{
				l.usernameInput,
				l.emailInput,
				l.passwordInput,
				l.passwordRepeatInput,
			}
			numSelections := len(inputs)


			// if submit register
			if key == tea.KeyEnter && l.index == numSelections {
				if l.passwordInput.Value() != l.passwordRepeatInput.Value() {
					return l, ErrorCmd(errors.New("passwords don't match"))
				}

				return NewRegisterModel(), register(l.usernameInput.Value(), l.emailInput.Value(), l.passwordInput.Value())
			}

			if key == tea.KeyUp || key == tea.KeyShiftTab {
				l.index--
			} else {
				l.index++
			}

			if l.index > numSelections {
				l.index = 0
			} else if l.index < 0 {
				l.index = numSelections
			}

			// Focus/Blur based on selection
			for i := 0; i < numSelections; i++ {
				if i == l.index {
					inputs[i].Focus()
					inputs[i].Prompt = focusedPrompt
					inputs[i].TextColor = focusedTextColor
				} else {
					inputs[i].Blur()
					inputs[i].Prompt = blurredPrompt
					inputs[i].TextColor = ""
				}
			}

			l.usernameInput = inputs[0]
			l.emailInput = inputs[1]
			l.passwordInput = inputs[2]
			l.passwordRepeatInput = inputs[3]

			if l.index == numSelections {
				l.submitButton = focusedSubmitButton
			} else {
				l.submitButton = blurredSubmitButton
			}
		case tea.KeyEsc:
			return NewRegisterModel(), ChangeViewCmd(ViewHome)
		}
	}

	l.usernameInput, cmd = l.usernameInput.Update(msg)
	cmds = append(cmds, cmd)

	l.emailInput, cmd = l.emailInput.Update(msg)
	cmds = append(cmds, cmd)

	l.passwordInput, cmd = l.passwordInput.Update(msg)
	cmds = append(cmds, cmd)

	l.passwordRepeatInput, cmd = l.passwordRepeatInput.Update(msg)
	cmds = append(cmds, cmd)

	return l, tea.Batch(cmds...)
}

func (l RegisterModel) View() (s string) {
	s = " Register\n\n"

	s += l.usernameInput.View() + "\n"
	s += l.emailInput.View() + "\n"
	s += l.passwordInput.View() + "\n"
	s += l.passwordRepeatInput.View() + "\n"
	s += "\n" + l.submitButton + "\n"

	return
}

