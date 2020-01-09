package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

// An interface for database connections, abstracting underlying DB access
type Connection interface {
	// should insert a user into the database and raise an error if it exists
	CreateUser(u User) error

	// Gets the user based on the username
	GetUserByName(username string) (User, error)
}

func InitDB(connectiontype string) Connection {
	// Database connection
	log.Printf("Starting up database ...")

	switch connectiontype {
	// in memory
	default:
		connection := SQLConnection{}

		var err error
		connection.db, err = gorm.Open("sqlite3", ":memory:")
		if err != nil {
			log.Fatal("Couldn't connect to the in memory sqlite3 database!")
		}

		// auto migrate schema
		connection.db.AutoMigrate(&User{})

		return connection
	}
}
