package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

type RepositoryCollection interface {
	UserRepository
}

// An interface for database connections, abstracting underlying DB access
type UserRepository interface {
	// should insert a user into the database and raise an error if it exists
	CreateUser(u User) error

	// Gets the user based on the username
	GetUserByName(username string) (User, error)

	// Counts the number of users in the database
	CountUsers() (int, error)

	// changes the fields of a user matching the username
	UpdateUser(user User) error

	// Gets users using specified start and end range
	GetUsers(start int, end int) ([]User, error)
}

// The connection types (wtb enums)
const (
	INMEMORY = "inmemory"
)

func InitDB(connectiontype string) RepositoryCollection {
	// Database connection
	log.Info("Starting up database ...")

	switch connectiontype {
	// in memory
	default:
		log.Debug("Using default in memory sqlite3 database")
		connection := SQLConnection{}

		var err error
		connection.db, err = gorm.Open("sqlite3", ":memory:")
		if err != nil {
			log.Fatal("Couldn't connect to the in memory sqlite3 database!")
		}

		// auto migrate schema
		connection.db.AutoMigrate(&userDAL{})

		return connection
	}
}
