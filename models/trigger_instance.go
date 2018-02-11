package models

import (
	"database/sql"
	"time"
)

type TriggerInstance struct {
	ID         int       `json:"id"`
	TriggerID  int       `json:"trigger_id"`
	InputData  string    `json:"input_data"`
	CreatedAt  time.Time `json:"created_at"`
}

func (t *TriggerInstance) GetTriggerInstance(db *sql.DB) error {
	return db.QueryRow("SELECT trigger_id, input_data FROM trigger_instances WHERE id=$1",
		t.ID).Scan(&t.TriggerID, &t.InputData)
}

func (t *TriggerInstance) CreateTriggerInstance(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO trigger_instances(trigger_id, input_data, created_at) VALUES($1, $2, $3) RETURNING id",
		t.TriggerID, t.InputData, t.CreatedAt).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}