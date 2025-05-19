package models

import (
	"fmt"

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

func (p Pair) String() string {
	return fmt.Sprintf("uuid: %s - %s/%s", p.UUID.String(), p.Numerator.Symbol, p.Denominator.Symbol)
}
