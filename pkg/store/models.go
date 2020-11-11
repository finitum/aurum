package store

import "github.com/finitum/aurum/pkg/models"

type DGraphUser struct {
	*models.User
	DType []string `json:"dgraph.type,omitempty"`
	Uid   string   `json:"uid,omitempty"`
}

func NewDGraphUser(user *models.User) *DGraphUser {
	return &DGraphUser{User: user, DType: []string{"User"}}

}

type DGraphApplication struct {
	models.Application
	DType []string `json:"dgraph.type,omitempty"`
	Uid   string   `json:"uid,omitempty"`
}

func NewDGraphApplication(application *models.Application) *DGraphApplication {
	return &DGraphApplication{Application: *application, DType: []string{"Application"}}
}
