package main

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/models"
	te "github.com/muesli/termenv"
	"strings"
)

var blurredBackButton = "[ " + te.String("Back").Foreground(color("240")).String() + " ]"
var focusedBackButton = "[ " + te.String("Back").Foreground(color(focusedTextColor)).String() + " ]"

var blurredNewGroupButton = "[ " + te.String("New group").Foreground(color("240")).String() + " ]"
var focusedNewGroupButton = "[ " + te.String("New group").Foreground(color(focusedTextColor)).String() + " ]"

var blurredDeleteGroupButton = "[ " + te.String("Delete group").Foreground(color("240")).String() + " ]"
var focusedDeleteGroupButton = "[ " + te.String("Delete group").Foreground(color(focusedTextColor)).String() + " ]"


type Button struct {
	name 	string
	blur    string
	focus   string
	focused bool
}

type GroupListModel struct {
	buttonIndex int
	listIndex   int
	buttons     []Button
	groups      []models.Group
	triedGroups bool
	newGroup    bool

	newGroupInput textinput.Model
}

func NewGroupListModel() GroupListModel {
	newGroupInput := textinput.NewModel()
	newGroupInput.Placeholder = "name"
	newGroupInput.Prompt = blurredPrompt

	return GroupListModel{
		buttonIndex: 0,
		listIndex: 0,
		buttons: []Button{
			{"back", blurredBackButton, focusedBackButton, true},
			{"new group", blurredNewGroupButton, focusedNewGroupButton, false},
			{"delete group", blurredDeleteGroupButton, focusedDeleteGroupButton, false},
		},
		groups: nil,
		triedGroups: false,
		newGroup: false,
		newGroupInput: newGroupInput,
	}
}

func (g GroupListModel) Update(msg tea.Msg) (GroupListModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)


	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.Type; key {
		case tea.KeyEnter:
			if g.newGroup {
				if g.newGroupInput.Value() == "" {
					return g, ErrorCmd(errors.New("group name can't be empty"))
				}
				g.newGroup = false
				g.buttons[g.buttonIndex].focused = true

				cmds = append(cmds, newGroup(g.newGroupInput.Value()))
			} else {
				switch g.buttons[g.buttonIndex].name {
				case "back":
					return NewGroupListModel(), ChangeViewCmd(ViewUser)
				case "new group":
					for i := 0; i < len(g.buttons); i++ {
						g.buttons[i].focused = false
					}

					g.newGroup = true
					g.newGroupInput.Reset()
					g.newGroupInput.Focus()
				case "delete group":
					if g.groups[g.listIndex].Name == aurum.AurumName {
						return g, ErrorCmd(errors.New("can't remove the aurum group"))
					}

					cmds = append(cmds, deleteGroup(g.groups[g.listIndex].Name))

					// Remove group from list
					copy(g.groups[g.listIndex:], g.groups[g.listIndex+1:])
					g.groups = g.groups[:len(g.groups)-1]
				}

			}

		case tea.KeyUp:
			if g.newGroup {
				break
			}
			g.listIndex --
			if g.listIndex < 0 {
				g.listIndex = len(g.groups) - 1
			}

		case tea.KeyDown:
			if g.newGroup {
				break
			}
			g.listIndex ++
			if g.listIndex > len(g.groups)-1 {
				g.listIndex = 0
			}
		case tea.KeyShiftTab, tea.KeyTab, tea.KeyLeft, tea.KeyRight:
			if g.newGroup {
				break
			}

			if key == tea.KeyLeft || key == tea.KeyShiftTab {
				g.buttonIndex--
			} else {
				g.buttonIndex++
			}

			if g.buttonIndex > len(g.buttons)-1 {
				g.buttonIndex = 0
			} else if g.buttonIndex < 0 {
				g.buttonIndex = len(g.buttons) - 1
			}

			// Focus/Blur based on selection
			for i := 0; i < len(g.buttons); i++ {
				if i == g.buttonIndex && !g.newGroup {
					g.buttons[i].focused = true
				} else {
					g.buttons[i].focused = false
				}
			}

		case tea.KeyEsc:
			if g.newGroup {
				g.newGroup = false
				g.buttons[g.buttonIndex].focused = true
			} else {
				return NewGroupListModel(), ChangeViewCmd(ViewUser)
			}
		}
	case GroupsMsg:
		g.groups = msg.groups
	}

	if !g.triedGroups {
		cmds = append(cmds, getGroups)
		g.triedGroups = true
	}


	g.newGroupInput, cmd = g.newGroupInput.Update(msg)
	cmds = append(cmds, cmd)

	return g, tea.Batch(cmds...)
}

func (g GroupListModel) View(width int) string {
	s := " Group list\n\n  "

	for _, i := range g.buttons {
		line := i.blur
		if i.focused {
			line = i.focus
		}

		s += line + " "
	}

	s += "\n"

	s += strings.Repeat("—", width) + "\n"

	for i, group := range g.groups {
		if i == g.listIndex {
			s += fmt.Sprintf(" * %s\n", te.String(group.Name).Foreground(color(focusedTextColor)).String())
		} else {
			s += fmt.Sprintf(" * %s\n", group.Name)
		}
	}

	if g.newGroup {
		s += "\n" + g.newGroupInput.View() + " Press enter to create \n\n"
	}

	return s
}
