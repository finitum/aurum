package main

import (
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ChangePasswordModel struct {
	index         int
	passwordInput textinput.Model
	passwordRepeatInput textinput.Model
	submitButton  string
}

func NewChangePasswordModel() ChangePasswordModel {
	password := textinput.NewModel()
	password.Placeholder = "password"
	password.Prompt = blurredPrompt
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.Focus()

	passwordrepeat := textinput.NewModel()
	passwordrepeat.Placeholder = "repeat password"
	passwordrepeat.Prompt = blurredPrompt
	passwordrepeat.EchoMode = textinput.EchoPassword
	passwordrepeat.EchoCharacter = '•'

	return ChangePasswordModel{
		index:         0,
		passwordInput: password,
		passwordRepeatInput: passwordrepeat,
		submitButton:  blurredSubmitButton,
	}
}

func (c ChangePasswordModel) Update(msg tea.Msg) (ChangePasswordModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.Type; key {
		case tea.KeyShiftTab, tea.KeyTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			inputs := []textinput.Model{
				c.passwordInput,
				c.passwordRepeatInput,
			}
			numSelections := len(inputs)


			// if submit register
			if key == tea.KeyEnter && c.index == numSelections {
				if c.passwordInput.Value() != c.passwordRepeatInput.Value() {
					return c, ErrorCmd(errors.New("passwords don't match"))
				}

				return NewChangePasswordModel(), changePassword(c.passwordInput.Value())
			}

			if key == tea.KeyUp || key == tea.KeyShiftTab {
				c.index--
			} else {
				c.index++
			}

			if c.index > numSelections {
				c.index = 0
			} else if c.index < 0 {
				c.index = numSelections
			}

			// Focus/Blur based on selection
			for i := 0; i < numSelections; i++ {
				if i == c.index {
					inputs[i].Focus()
					inputs[i].Prompt = focusedPrompt
					inputs[i].TextColor = focusedTextColor
				} else {
					inputs[i].Blur()
					inputs[i].Prompt = blurredPrompt
					inputs[i].TextColor = ""
				}
			}

			c.passwordInput = inputs[0]
			c.passwordRepeatInput = inputs[1]

			if c.index == numSelections {
				c.submitButton = focusedSubmitButton
			} else {
				c.submitButton = blurredSubmitButton
			}
		case tea.KeyEsc:
			return NewChangePasswordModel(), ChangeViewCmd(ViewUser)
		}
	}

	c.passwordInput, cmd = c.passwordInput.Update(msg)
	cmds = append(cmds, cmd)

	c.passwordRepeatInput, cmd = c.passwordRepeatInput.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c ChangePasswordModel) View() (s string) {
	s = " Change password\n\n"

	s += c.passwordInput.View() + "\n"
	s += c.passwordRepeatInput.View() + "\n"
	s += "\n" + c.submitButton + "\n"

	return
}

