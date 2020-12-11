package main

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/models"
	te "github.com/muesli/termenv"
	"strings"
)

type EditUserGroupModel struct {
	buttonIndex int
	listIndex   int
	buttons     []Button
	userGroups  []models.GroupWithRole
	allGroups   []models.Group
	triedGroups bool
	triedAllGroups bool
	newGroup    bool

	user          *models.User
}

var blurredAdminButton = "[ " + te.String("Toggle admin").Foreground(color("240")).String() + " ]"
var focusedAdminButton = "[ " + te.String("Toggle admin").Foreground(color(focusedTextColor)).String() + " ]"


func NewEditUserGroupModel() EditUserGroupModel {
	return EditUserGroupModel{
		buttonIndex: 0,
		listIndex:   0,
		buttons: []Button{
			{"back", blurredBackButton, focusedBackButton, true},
			{"new group", blurredNewGroupButton, focusedNewGroupButton, false},
			{"delete group", blurredDeleteGroupButton, focusedDeleteGroupButton, false},
			{"toggle admin", blurredAdminButton, focusedAdminButton, false},
		},
		userGroups:  nil,
		triedGroups: false,
		newGroup:    false,
		user:        nil,
	}
}

func (g EditUserGroupModel) Update(msg tea.Msg) (EditUserGroupModel, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.Type; key {
		case tea.KeyEnter:
			if g.newGroup {
				g.newGroup = false
				g.buttons[g.buttonIndex].focused = true
				cmds = append(cmds, addUserToGroup(g.user.Username, g.allGroups[g.listIndex].Name))

				g.userGroups = append(g.userGroups, models.GroupWithRole{
					Group: g.allGroups[g.listIndex],
					Role:  models.RoleUser,
				})
				g.listIndex = 0
			} else {
				switch g.buttons[g.buttonIndex].name {
				case "back":
					return NewEditUserGroupModel(), ChangeViewCmd(ViewUserList)
				case "new group":
					for i := 0; i < len(g.buttons); i++ {
						g.buttons[i].focused = false
					}
					g.newGroup = true
					g.listIndex = 0
				case "toggle admin":
					if g.userGroups[g.listIndex].Role == models.RoleAdmin {
						cmds = append(cmds, setAccess(g.user.Username, g.userGroups[g.listIndex].Name, models.RoleUser))
						g.userGroups[g.listIndex].Role = models.RoleUser
					} else {
						cmds = append(cmds, setAccess(g.user.Username, g.userGroups[g.listIndex].Name, models.RoleAdmin))
						g.userGroups[g.listIndex].Role = models.RoleAdmin
					}

				case "delete group":
					if g.userGroups[g.listIndex].Name == aurum.AurumName {
						return g, ErrorCmd(errors.New("can't remove user from the Aurum group"))
					}

					cmds = append(cmds, removeUserFromGroup(g.user.Username, g.userGroups[g.listIndex].Name))

					// Remove group from list
					copy(g.userGroups[g.listIndex:], g.userGroups[g.listIndex+1:])
					g.userGroups = g.userGroups[:len(g.userGroups)-1]
				}
			}

		case tea.KeyUp:
			g.listIndex--
			if g.listIndex < 0 {
				if g.newGroup {
					g.listIndex = len(g.allGroups) - 1
				} else {
					g.listIndex = len(g.userGroups) - 1
				}
			}

		case tea.KeyDown:
			g.listIndex++
			if g.newGroup {
				if g.listIndex > len(g.allGroups)-1 {
					g.listIndex = 0
				}
			} else {
				if g.listIndex > len(g.userGroups)-1 {
					g.listIndex = 0
				}
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
				return NewEditUserGroupModel(), ChangeViewCmd(ViewUserList)
			}
		}
	case UserMsg:
		g.user = &msg.user
		g.triedGroups = false
	case UserGroupsMsg:
		g.userGroups = msg.groups
	case GroupsMsg:
		g.allGroups = msg.groups
	}

	if !g.triedGroups && g.user != nil {
		cmds = append(cmds, getGroupsForUser(g.user.Username, 0))
		g.triedGroups = true
	}

	if !g.triedAllGroups {
		cmds = append(cmds, getGroups)
		g.triedAllGroups = true
	}

	return g, tea.Batch(cmds...)
}

func (g EditUserGroupModel) View(width int) string {
	s := " Edit Groups \n\n"

	if g.user != nil {
		s += " Actions here apply for user: " + g.user.Username + "\n\n "
	}

	for _, i := range g.buttons {
		line := i.blur
		if i.focused {
			line = i.focus
		}

		s += line + " "
	}

	s += "\n"

	s += strings.Repeat("—", width) + "\n"

	if g.newGroup {
		s += "Choose a group: \n"
		for i, group := range g.allGroups {
			if i == g.listIndex {
				s += fmt.Sprintf(" * %s\n", te.String(group.Name).Foreground(color(focusedTextColor)).String())
			} else {
				s += fmt.Sprintf(" * %s\n", group.Name)
			}
		}
	} else {
		for i, group := range g.userGroups {
			if i == g.listIndex {
				s += fmt.Sprintf(" * %s  %s\n", te.String(fmt.Sprintf("%-20s", group.Name)).Foreground(color(focusedTextColor)).String(), group.Role.String())
			} else {
				s += fmt.Sprintf(" * %s  %s\n", fmt.Sprintf("%-20s", group.Name), group.Role.String())
			}
		}
	}

	return s
}

func (g EditUserGroupModel) Init(params []interface{}) (EditUserGroupModel, tea.Cmd) {
	if len(params) != 1 {
		return g, ErrorCmd(errors.New("failed to load user for editing"))
	}

	var ok bool
	g.user, ok = params[0].(*models.User)
	g.triedGroups = false
	g.triedAllGroups = false
	if !ok {
		return g, ErrorCmd(errors.New("failed to load user for editing"))
	}

	return g, nil
}
