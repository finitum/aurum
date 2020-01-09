package db

import (
	"aurum/util"
	"github.com/jinzhu/gorm"
)

type SQLConnection struct {
	db *gorm.DB
}

func (conn SQLConnection) CreateUser(u User) error {
	pass, err := util.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = pass

	// Will error if user already exists
	d := conn.db.Create(&u)
	if d.Error != nil {
		return d.Error
	}

	return nil
}

func (conn SQLConnection) GetUserByName(username string) (User, error) {
	var u User

	if d := conn.db.Where(&User{
		Username: username,
	}).First(&u); d.Error != nil {
		return u, d.Error
	}

	return u, nil
}
