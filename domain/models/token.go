package models

import "github.com/google/uuid"

type Token struct {
	UUID      uuid.UUID `db:"uuid"`
	Name      string    `db:"name"`
	Symbol    string    `db:"symbol"`
	NetworkID string    `db:"network_id"`
	Network   string    `db:"network"`
}
