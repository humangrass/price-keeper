package repository

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/humangrass/gommon/database"
	"github.com/humangrass/price-keeper/domain/entities"
	"github.com/humangrass/price-keeper/domain/models"
	"github.com/humangrass/price-keeper/pgk/x/xtype"
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

type pairScan struct {
	UUID      uuid.UUID      `db:"uuid"`
	Timeframe xtype.Interval `db:"timeframe"`
	IsActive  bool           `db:"is_active"`

	NumUUID      uuid.UUID `db:"numerator_uuid"`
	NumName      string    `db:"numerator_name"`
	NumSymbol    string    `db:"numerator_symbol"`
	NumNetworkID string    `db:"numerator_network_id"`
	NumNetwork   string    `db:"numerator_network"`

	DenUUID      uuid.UUID `db:"denominator_uuid"`
	DenName      string    `db:"denominator_name"`
	DenSymbol    string    `db:"denominator_symbol"`
	DenNetworkID string    `db:"denominator_network_id"`
	DenNetwork   string    `db:"denominator_network"`
}

func (r *PairsRepository) GetByParams(ctx context.Context, params entities.RequestParams) ([]models.Pair, int, error) {
	var scans []pairScan
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
			goqu.I("p.uuid"),
			goqu.I("p.timeframe"),
			goqu.I("p.is_active"),

			goqu.I("t1.uuid").As("numerator_uuid"),
			goqu.I("t1.name").As("numerator_name"),
			goqu.I("t1.symbol").As("numerator_symbol"),
			goqu.I("t1.network_id").As("numerator_network_id"),
			goqu.I("t1.network").As("numerator_network"),

			goqu.I("t2.uuid").As("denominator_uuid"),
			goqu.I("t2.name").As("denominator_name"),
			goqu.I("t2.symbol").As("denominator_symbol"),
			goqu.I("t2.network_id").As("denominator_network_id"),
			goqu.I("t2.network").As("denominator_network"),
		).
		From(goqu.T(PairsTableName).As("p")).
		Join(
			goqu.T(TokensTableName).As("t1"),
			goqu.On(goqu.I("p.numerator").Eq(goqu.I("t1.uuid"))),
		).
		Join(
			goqu.T(TokensTableName).As("t2"),
			goqu.On(goqu.I("p.denominator").Eq(goqu.I("t2.uuid"))),
		).
		Where(goqu.I("t1.network").Eq(goqu.I("t2.network"))).
		Limit(uint(params.Limit)).
		Offset(uint(params.Offset))

	if params.OrderBy == entities.OrderByDesc {
		query = query.Order(goqu.I("p.timeframe").Desc())
	} else {
		query = query.Order(goqu.I("p.timeframe").Asc())
	}
	err = query.ScanStructsContext(ctx, &scans)
	pairs := make([]models.Pair, len(scans))
	for i, s := range scans {
		pairs[i] = models.Pair{
			UUID:      s.UUID,
			Timeframe: s.Timeframe,
			IsActive:  s.IsActive,
			Numerator: models.Token{
				UUID:      s.NumUUID,
				Name:      s.NumName,
				Symbol:    s.NumSymbol,
				NetworkID: s.NumNetworkID,
				Network:   s.NumNetwork,
			},
			Denominator: models.Token{
				UUID:      s.DenUUID,
				Name:      s.DenName,
				Symbol:    s.DenSymbol,
				NetworkID: s.DenNetworkID,
				Network:   s.DenNetwork,
			},
		}
	}

	return pairs, total, err
}

func (r *PairsRepository) Create(ctx context.Context, pair *models.Pair) error {
	if pair.UUID == uuid.Nil {
		pair.UUID = uuid.New()
	}

	_, err := r.pool.Builder().
		Insert(PairsTableName).
		Rows(goqu.Record{
			"uuid":        pair.UUID,
			"numerator":   pair.Numerator.UUID,
			"denominator": pair.Denominator.UUID,
			"timeframe":   pair.Timeframe,
			"is_active":   pair.IsActive,
		}).Executor().
		ExecContext(ctx)

	return err
}

func (r *PairsRepository) Update(ctx context.Context, pair *models.Pair) error {
	if pair.UUID == uuid.Nil {
		return fmt.Errorf("pair UUID is required for update")
	}

	updateData := goqu.Record{
		"numerator":   pair.Numerator.UUID,
		"denominator": pair.Denominator.UUID,
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
