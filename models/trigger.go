package models

import (
	"database/sql"
)

type Trigger struct {
	ID       int    `json:"id"`
	DeviceID int    `json:"device_id"`
	Key      string `json:"key"`
}

func (t *Trigger) GetTrigger(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM triggers WHERE device_id=$1 AND key=$2",
		t.DeviceID, t.Key).Scan(&t.ID)
}