package db

import (
	"errors"
	"github.com/finitum/aurum/pkg/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

var ErrExists = errors.New("user already exists")

type RepositoryCollection interface {
	UserRepository
}

// An interface for database connections, abstracting underlying DB access
type UserRepository interface {
	// should insert a user into the database and raise ErrExists if the user already exists
	CreateUser(u models.User) error

	// Gets the user based on the username
	GetUserByName(username string) (models.User, error)

	// Counts the number of users in the database
	CountUsers() (int, error)

	// changes the fields of a user matching the username
	UpdateUser(user models.User) error

	// Gets users using specified start and end range
	GetUsers(start int, end int) ([]models.User, error)
}

type DatabaseType int

const (
	InMemory DatabaseType = iota
)

func InitDB(db DatabaseType) RepositoryCollection {
	// Database connection
	log.Info("Starting up database ...")

	switch db {
	// in memory
	case InMemory:
		fallthrough
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
