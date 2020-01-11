package db

import (
	"aurum/hash"
	"github.com/jinzhu/gorm"
)

type SQLConnection struct {
	db *gorm.DB
}

func (conn SQLConnection) CreateUser(u User) error {
	pass, err := hash.HashPassword(u.Password)
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

func (conn SQLConnection) CountUsers() (int, error) {

	var i int
	if d := conn.db.Model(&User{}).Count(&i); d.Error != nil {
		return 0, d.Error
	}

	return i, nil
}

func (conn SQLConnection) UpdateUser(user User) error {

	// Will error if user already exists
	d := conn.db.Save(&user)
	if d.Error != nil {
		return d.Error
	}

	return nil
}

func (conn SQLConnection) GetUsers(start int, end int) ([]User, error) {
	var users []User
	d := conn.db.Model(&User{}).Offset(start).Limit(end - start).Find(&users)
	if d.Error != nil {
		return nil, d.Error
	}

	return users, nil
}
