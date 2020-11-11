package dgraph

import "github.com/finitum/aurum/pkg/models"

type User struct {
	*models.User

	DType []string `json:"dgraph.type,omitempty"`
	Uid   string   `json:"uid,omitempty"`

	Applications []Application `json:"applications,omitempty"`
}

type Application struct {
	models.Application

	Role models.Role `json:"applications|role,omitempty"`

	DType []string `json:"dgraph.type,omitempty"`
	Uid   string   `json:"uid,omitempty"`
}

func NewDGraphUser(user *models.User) *User {
	return &User{User: user, DType: []string{"User"}}
}



func NewDGraphApplication(application *models.Application) *Application {
	return &Application{Application: *application, DType: []string{"Application"}}
}
