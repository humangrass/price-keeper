package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/humangrass/gommon/database"
	"github.com/humangrass/price-keeper/domain/models"
	"github.com/humangrass/price-keeper/pgk/x/xtype"
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
	FindBySymbols(ctx context.Context, numerator, denominator string) ([]models.Price, error)
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

type priceScan struct {
	UUID  uuid.UUID `db:"uuid"`
	TS    time.Time `db:"ts"`
	Price float64   `db:"price"`

	PairUUID  uuid.UUID      `db:"pair_uuid"`
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

func (r *PricesRepository) priceQuery() *goqu.SelectDataset {
	return r.pool.Builder().
		Select(
			goqu.I("pr.uuid"),
			goqu.I("pr.ts"),
			goqu.I("pr.price"),

			goqu.I("p.uuid").As("pair_uuid"),
			goqu.I("p.timeframe").As("timeframe"),
			goqu.I("p.is_active").As("is_active"),

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
		From(goqu.T("prices").As("pr")).
		Join(
			goqu.T("pairs").As("p"),
			goqu.On(goqu.I("pr.pair_uuid").Eq(goqu.I("p.uuid"))),
		).
		Join(
			goqu.T("tokens").As("t1"),
			goqu.On(goqu.I("p.numerator").Eq(goqu.I("t1.uuid"))),
		).
		Join(
			goqu.T("tokens").As("t2"),
			goqu.On(goqu.I("p.denominator").Eq(goqu.I("t2.uuid"))),
		)
}

func (r *PricesRepository) FindBySymbols(ctx context.Context, numerator, denominator string) ([]models.Price, error) {
	var scans []priceScan

	query := r.priceQuery().
		Where(
			goqu.Func("LOWER", goqu.I("t1.symbol")).Eq(strings.ToLower(numerator)),
			goqu.Func("LOWER", goqu.I("t2.symbol")).Eq(strings.ToLower(denominator)),
			goqu.I("t1.network").Eq(goqu.I("t2.network")),
		).
		Order(goqu.I("pr.ts").Desc())

	err := query.ScanStructsContext(ctx, &scans)
	if err != nil {
		return nil, fmt.Errorf("failed to scan prices: %w", err)
	}

	prices := make([]models.Price, len(scans))
	for i, s := range scans {
		prices[i] = models.Price{
			UUID:  s.UUID,
			TS:    s.TS,
			Price: s.Price,
			Pair: &models.Pair{
				UUID:      s.PairUUID,
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
			},
		}
	}

	return prices, nil
}
