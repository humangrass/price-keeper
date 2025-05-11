package models

import (
	"time"

	"github.com/google/uuid"
)

type Pair struct {
	UUID        uuid.UUID     `db:"uuid"`
	Numerator   uuid.UUID     `db:"numerator"`
	Denominator uuid.UUID     `db:"denominator"`
	Timeframe   time.Duration `db:"timeframe"`
	IsActive    bool          `db:"is_activce"`
}
