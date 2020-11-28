package dgraph

import "github.com/finitum/aurum/pkg/models"

type User struct {
	models.User

	DType []string `json:"dgraph.type,omitempty"`
	Uid   string   `json:"uid,omitempty"`

	Groups []Group `json:"groups,omitempty"`
}

type Group struct {
	models.Group

	Role models.Role `json:"groups|role,omitempty"`

	DType []string `json:"dgraph.type,omitempty"`
	Uid   string   `json:"uid,omitempty"`
}

func NewDGraphUser(user models.User) *User {
	return &User{User: user, DType: []string{"User"}}
}

func NewDGraphGroup(group models.Group) *Group {
	return &Group{Group: group, DType: []string{"Group"}}
}
