package db

import (
	"github.com/finitum/aurum/pkg/hash"
	"github.com/finitum/aurum/pkg/models"
	"github.com/jinzhu/gorm"
	"strings"
)

// The database model for a Gorm user
type userDAL struct {
	gorm.Model  `json:"-"`
	models.User `gorm:"embedded"`
}

type SQLConnection struct {
	db *gorm.DB
}

func (conn SQLConnection) CreateUser(u models.User) error {
	pass, err := hash.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = pass

	// Will error if user already exists
	d := conn.db.Create(&userDAL{User: u})

	if d.Error != nil {
		if strings.HasPrefix(d.Error.Error(), "UNIQUE constraint failed:") {
			return ErrExists
		}
		return d.Error
	}

	return nil
}

func (conn SQLConnection) GetUserByName(username string) (models.User, error) {
	var u = &userDAL{}
	u.Username = username

	if d := conn.db.Where(u).First(&u); d.Error != nil {
		return models.User{}, d.Error
	}

	return u.User, nil
}

func (conn SQLConnection) CountUsers() (int, error) {

	var i int
	if d := conn.db.Model(&userDAL{}).Count(&i); d.Error != nil {
		return 0, d.Error
	}

	return i, nil
}

func (conn SQLConnection) UpdateUser(user models.User) error {
	var dbuser = &userDAL{}
	dbuser.Username = user.Username

	if d := conn.db.Where(dbuser).First(&dbuser); d.Error != nil {
		return d.Error
	}

	dbuser.User = user

	d := conn.db.Save(&dbuser)
	if d.Error != nil {
		return d.Error
	}

	return nil
}

func (conn SQLConnection) GetUsers(start int, end int) ([]models.User, error) {
	var users []userDAL
	d := conn.db.Model(&userDAL{}).Offset(start).Limit(end - start).Find(&users)
	if d.Error != nil {
		return nil, d.Error
	}

	var us []models.User

	for _, element := range users {
		us = append(us, element.User)
	}

	return us, nil
}
