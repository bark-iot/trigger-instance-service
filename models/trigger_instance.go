package models

import (
	"database/sql"
)

type TriggerInstance struct {
	ID         int    `json:"id"`
	TriggerID  int    `json:"trigger_id"`
	InputData  string `json:"input_data"`
}

func (t *TriggerInstance) GetTriggerInstance(db *sql.DB) error {
	return db.QueryRow("SELECT trigger_id, input_data FROM trigger_instances WHERE id=$1",
		t.ID).Scan(&t.TriggerID, &t.InputData)
}

func (t *TriggerInstance) CreateTriggerInstance(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO trigger_instances(trigger_id, input_data) VALUES($1, $2) RETURNING id",
		t.TriggerID, t.InputData).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}