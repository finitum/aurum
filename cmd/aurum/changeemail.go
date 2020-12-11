package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/pkg/models"
)

type ChangeEmailModel struct {
	index        int
	emailInput   textinput.Model
	user         *models.User
	submitButton string
}

func NewChangeEmailModel() ChangeEmailModel {
	email := textinput.NewModel()
	email.Placeholder = "email"
	email.Prompt = blurredPrompt
	email.Focus()

	return ChangeEmailModel{
		index:        0,
		emailInput:   email,
		submitButton: blurredSubmitButton,
	}
}

func (c ChangeEmailModel) Update(msg tea.Msg) (ChangeEmailModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.Type; key {
		case tea.KeyShiftTab, tea.KeyTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			inputs := []textinput.Model{
				c.emailInput,
			}
			numSelections := len(inputs)

			// if submit register
			if key == tea.KeyEnter && c.index == numSelections {
				return NewChangeEmailModel(), changeEmail(c.emailInput.Value())
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

			c.emailInput = inputs[0]

			if c.index == numSelections {
				c.submitButton = focusedSubmitButton
			} else {
				c.submitButton = blurredSubmitButton
			}
		case tea.KeyEsc:
			return NewChangeEmailModel(), ChangeViewCmd(ViewUser)
		}
	case UserMsg:
		c.user = &msg.user
	}

	c.emailInput, cmd = c.emailInput.Update(msg)
	cmds = append(cmds, cmd)

	if c.user == nil {
		cmds = append(cmds, getUser)
	}


	return c, tea.Batch(cmds...)
}

func (c ChangeEmailModel) View() (s string) {
	s = " Change email address\n\n"

	if c.user != nil {
		s += "Previous email: " + c.user.Email + "\n\n"
	}
	s += c.emailInput.View() + "\n"
	s += "\n" + c.submitButton + "\n"

	return
}
