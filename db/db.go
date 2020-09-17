package db

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Database is the core database interface
type Database interface {
	AddEvent(resourceID string, resourceType, eventType int32, data []byte)
}

// TinkDB implements the Database interface
type TinkDB struct {
	instance *sql.DB
}

// Connect returns a connection to postgres database
func Connect() Database {
	// const connInfo = "dbname=tinkerbell user=tinkerbell password=tinkerbell sslmode=disable"
	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return &TinkDB{instance: db}
}
