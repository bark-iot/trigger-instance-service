package models

import (
	"database/sql"
)

type Device struct {
	ID     int    `json:"id"`
	Token string  `json:"token"`
}

func (d *Device) GetDevice(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM devices WHERE token=$1",
		d.Token).Scan(&d.ID)
}