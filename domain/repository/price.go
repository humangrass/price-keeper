package repository

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/humangrass/gommon/database"
	"github.com/humangrass/price-keeper/domain/models"
)

const PricesTableName = "prices"

type PricesRepository struct {
	BaseRepository
}

func NewPricesRepository(pool database.Pool) *PricesRepository {
	return &PricesRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

type PriceRepo interface {
	Create(ctx context.Context, model models.Price) error
}

func (r *PricesRepository) Create(ctx context.Context, model models.Price) error {
	if model.UUID == uuid.Nil {
		model.UUID = uuid.New()
	}

	_, err := r.pool.Builder().
		Insert(PricesTableName).
		Rows(goqu.Record{
			"uuid":      model.UUID,
			"ts":        time.Now(),
			"price":     model.Price,
			"pair_uuid": model.Pair.UUID,
		}).Executor().
		ExecContext(ctx)

	return err
}
