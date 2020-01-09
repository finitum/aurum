package db

import "github.com/jinzhu/gorm"

// The various roles a user can be
const (
	rUSER  int = 0
	rADMIN int = 1
)

// The database model for a user
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username, omitempty"`
	Password string `gorm:"not null" json:"password,omitempty"`
	Email    string `gorm:"not null" json:"email,omitempty"`
	Role     int    `gorm:"not null" sql:"DEFAULT:0" json:"role,omitempty"`
	Blocked  bool 	`gorm:"not null" sql:"DEFAULT:false" json:"blocked,omitempty"`
}

