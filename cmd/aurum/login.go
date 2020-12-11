package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	te "github.com/muesli/termenv"
)

const focusedTextColor = "205"

var (
	focusedPrompt       = te.String("> ").Foreground(color(focusedTextColor)).String()
	blurredPrompt       = "> "
	focusedSubmitButton = "[ " + te.String("Submit").Foreground(color(focusedTextColor)).String() + " ]"
	blurredSubmitButton = "[ " + te.String("Submit").Foreground(color("240")).String() + " ]"
)


type LoginModel struct {
	index         int
	usernameInput textinput.Model
	passwordInput textinput.Model
	submitButton  string
}

func NewLoginModel() LoginModel {
	name := textinput.NewModel()
	name.Placeholder = "username"
	name.Focus()
	name.Prompt = focusedPrompt
	name.TextColor = focusedTextColor

	password := textinput.NewModel()
	password.Placeholder = "password"
	password.Prompt = blurredPrompt
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	return LoginModel{
		index:         0,
		usernameInput: name,
		passwordInput: password,
		submitButton:  blurredSubmitButton,
	}
}

func (l LoginModel) Update(msg tea.Msg) (LoginModel, tea.Cmd) {
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
				l.passwordInput,
			}
			numSelections := len(inputs)

			// if submit login
			if key == tea.KeyEnter && l.index == numSelections {
				return NewLoginModel(), login(l.usernameInput.Value(), l.passwordInput.Value())
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
			l.passwordInput = inputs[1]

			if l.index == numSelections {
				l.submitButton = focusedSubmitButton
			} else {
				l.submitButton = blurredSubmitButton
			}
		case tea.KeyEsc:
			return NewLoginModel(), ChangeViewCmd(ViewHome)
		}
	}

	l.usernameInput, cmd = l.usernameInput.Update(msg)
	cmds = append(cmds, cmd)

	l.passwordInput, cmd = l.passwordInput.Update(msg)
	cmds = append(cmds, cmd)

	return l, tea.Batch(cmds...)
}

func (l LoginModel) View() (s string) {
	s = " Login\n\n"

	s += l.usernameInput.View() + "\n"
	s += l.passwordInput.View() + "\n\n"
	s += "\n" + l.submitButton + "\n"

	return
}

