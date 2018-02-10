package models

import (
	"database/sql"
)

type House struct {
	ID     int    `json:"id"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

func (h *House) GetHouse(db *sql.DB) error {
	return db.QueryRow("SELECT id, secret FROM houses WHERE key=$1",
		h.Key).Scan(&h.ID, &h.Secret)
}