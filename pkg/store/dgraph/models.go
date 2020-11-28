package dgraph

import "github.com/finitum/aurum/pkg/models"

type User struct {
	models.User

	DType []string `json:"dgraph.type,omitempty"`
	Uid   string   `json:"uid,omitempty"`
}

func NewDGraphUser(user models.User) *User {
	return &User{User: user, DType: []string{"User"}}
}
