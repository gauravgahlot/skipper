package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// AddEvent writes a new workflow event into the database
func (d TinkDB) AddEvent(resourceID string, resourceType, eventType int32, data []byte) {
	tx, err := d.instance.BeginTx(context.TODO(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error(err)
		return
	}

	_, err = tx.Exec(`
	INSERT INTO
		events (id, resource_id, resource_type, event_type, data, created_at)
	VALUES
		($1, $2, $3, $4, $5, $6)
	`, uuid.New(), resourceID, resourceType, eventType, data, time.Now())

	if err != nil {
		log.Error(err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return
	}
}
