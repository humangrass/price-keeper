package models

import (
	"time"

	"github.com/google/uuid"
)

type Price struct {
	UUID  uuid.UUID `db:"uuid"`
	TS    time.Time `db:"ts"`
	Price float64   `db:"price"`
	Pair  *Pair     `db:"pair_uuid"`
}
