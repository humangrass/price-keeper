package models

import (
	"github.com/google/uuid"
	"github.com/humangrass/price-keeper/pgk/x/xtype"
)

type Pair struct {
	UUID        uuid.UUID      `db:"uuid"`
	Numerator   Token          `db:"numerator"`
	Denominator Token          `db:"denominator"`
	Timeframe   xtype.Interval `db:"timeframe"`
	IsActive    bool           `db:"is_active"`
}
