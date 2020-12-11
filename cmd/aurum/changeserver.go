package main

import (
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	te "github.com/muesli/termenv"
)

type ChangeServerModel struct {
	index     int
	prevView  View
	newClient textinput.Model
}

func NewChangeServerModel() ChangeServerModel {

	newClient := textinput.NewModel()
	newClient.Placeholder = "add url"
	newClient.TextColor = focusedTextColor
	newClient.Prompt = "  "

	return ChangeServerModel{
		index:     0,
		newClient: newClient,
	}
}

func (c ChangeServerModel) Update(msg tea.Msg) (ChangeServerModel, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.Type; key {
		case tea.KeyEnter:
			if c.index == clientManager.NumClients() {
				host := c.newClient.Value()
				c.newClient.Reset()

				return c, func() tea.Msg {
					err := clientManager.AddClient(host)
					if err != nil {
						return ErrorMsg{err}
					}

					err = clientManager.SetActiveClient(clientManager.NumClients() - 1)
					if err != nil {
						return ErrorMsg{err}
					}

					return nil
				}

			} else {
				err := clientManager.SetActiveClient(c.index)
				if err != nil {
					return c, ErrorCmd(err)
				}

				return NewChangeServerModel(), ChangeViewCmd(c.prevView)
			}

		case tea.KeyUp:
			c.index--
			if c.index < 0 {
				c.index = clientManager.NumClients()
			}

		case tea.KeyDown:
			c.index++
			if c.index > clientManager.NumClients() {
				c.index = 0
			}
		case tea.KeyEsc:
			return NewChangeServerModel(), ChangeViewCmd(c.prevView)
		}
	}

	if c.index == clientManager.NumClients() {
		c.newClient.Focus()
	} else {
		c.newClient.Blur()
	}

	c.newClient, cmd = c.newClient.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c ChangeServerModel) View() (s string) {
	s = " Servers \n\n"

	active, err := clientManager.GetActiveClient()
	var activeUrl string
	if err != nil {
		activeUrl = ""
	} else {
		activeUrl = active.GetUrl()
	}

	for i := 0; i < clientManager.NumClients(); i++ {
		client, _ := clientManager.GetClient(i)
		url := client.GetUrl()

		if i == c.index {
			url = te.String(url).Foreground(color(focusedTextColor)).String()
		}

		if client.GetUrl() == activeUrl {
			s += "> " + url + "\n"
		} else {
			s += "  " + url + "\n"
		}
	}

	s += c.newClient.View()

	return
}

func (c ChangeServerModel) Init(params []interface{}) (ChangeServerModel, tea.Cmd) {
	if len(params) != 1 {
		return c, ErrorCmd(errors.New("failed to open change server view"))
	}

	var ok bool
	c.prevView, ok = params[0].(View)
	if !ok {
		return c, ErrorCmd(errors.New("failed to open change server view"))
	}

	c.index = clientManager.GetActiveClientIndex()

	return c, nil
}
