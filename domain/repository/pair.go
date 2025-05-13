package repository

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/humangrass/gommon/database"
	"github.com/humangrass/price-keeper/domain/entities"
	"github.com/humangrass/price-keeper/domain/models"
)

const PairsTableName = "pairs"

type PairsRepository struct {
	BaseRepository
}

func NewPairsRepository(pool database.Pool) *PairsRepository {
	return &PairsRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

type PairsRepo interface {
	GetByParams(ctx context.Context, params entities.RequestParams) ([]models.Pair, int, error)
	Create(ctx context.Context, pair *models.Pair) error
	Update(ctx context.Context, pair *models.Pair) error
}

func (r *PairsRepository) GetByParams(ctx context.Context, params entities.RequestParams) ([]models.Pair, int, error) {
	var pairs []models.Pair
	var total int

	_, err := r.pool.Builder().
		Select(goqu.COUNT("uuid")).
		From(PairsTableName).
		ScanVal(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pairs: %w", err)
	}

	query := r.pool.Builder().
		Select(
			goqu.C("uuid"),
			goqu.C("numerator"),
			goqu.C("denominator"),
			goqu.C("timeframe"),
			goqu.C("is_active"),
		).
		From(PairsTableName).
		Limit(uint(params.Limit)).
		Offset(uint(params.Offset))

	switch params.OrderBy {
	case entities.OrderByAsc:
		query = query.Order(goqu.C("timeframe").Asc())
	case entities.OrderByDesc:
		query = query.Order(goqu.C("timeframe").Desc())
	}
	err = query.ScanStructsContext(ctx, &pairs)
	return pairs, total, err
}

func (r *PairsRepository) Create(ctx context.Context, pair *models.Pair) error {
	if pair.UUID == uuid.Nil {
		pair.UUID = uuid.New()
	}

	_, err := r.pool.Builder().
		Insert(PairsTableName).
		Rows(pair).
		Executor().ExecContext(ctx)

	return err
}

func (r *PairsRepository) Update(ctx context.Context, pair *models.Pair) error {
	if pair.UUID == uuid.Nil {
		return fmt.Errorf("pair UUID is required for update")
	}

	updateData := goqu.Record{
		"numerator":   pair.Numerator,
		"denominator": pair.Denominator,
		"timeframe":   pair.Timeframe,
		"is_active":   pair.IsActive,
	}

	query, args, err := r.pool.Builder().
		Update(PairsTableName).
		Set(updateData).
		Where(goqu.C("uuid").Eq(pair.UUID)).
		ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.pool.Builder().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update pair: %w", err)
	}

	return nil
}
