package repository

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/humangrass/gommon/database"
	"github.com/humangrass/price-keeper/domain/entities"
	"github.com/humangrass/price-keeper/domain/models"
	"github.com/humangrass/price-keeper/pgk/xerror"
)

const TokensTableName = "tokens"

type TokensRepository struct {
	BaseRepository
}

func NewTokensRepository(pool database.Pool) *TokensRepository {
	return &TokensRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

type TokenRepo interface {
	GetByParams(ctx context.Context, params entities.RequestParams) ([]models.Token, int, error)
	Create(ctx context.Context, token *models.Token) error
	GetTokenBySymbol(ctx context.Context, symbol string) (models.Token, error)
}

type TokenSelect struct {
	models.Token
}

func (r *TokensRepository) GetByParams(ctx context.Context, params entities.RequestParams) ([]models.Token, int, error) {
	var tokens []models.Token
	var total int

	_, err := r.pool.Builder().
		Select(goqu.COUNT("uuid")).
		From(TokensTableName).
		ScanVal(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tokens: %w", err)
	}

	query := r.pool.Builder().
		Select(
			goqu.C("uuid"),
			goqu.C("name"),
			goqu.C("symbol"),
			goqu.C("network_id"),
			goqu.C("network"),
		).
		From(TokensTableName).
		Limit(uint(params.Limit)).
		Offset(uint(params.Offset))

	switch params.OrderBy {
	case entities.OrderByAsc:
		query = query.Order(goqu.C("name").Asc())
	case entities.OrderByDesc:
		query = query.Order(goqu.C("name").Desc())
	}

	err = query.ScanStructsContext(ctx, &tokens)
	return tokens, total, err
}

func (r *TokensRepository) Create(ctx context.Context, token *models.Token) error {
	if token.UUID == uuid.Nil {
		token.UUID = uuid.New()
	}

	_, err := r.pool.Builder().
		Insert(TokensTableName).
		Rows(token).
		Executor().
		ExecContext(ctx)

	return err
}

func (r *TokensRepository) GetTokenBySymbol(ctx context.Context, symbol string) (models.Token, error) {
	var token models.Token

	found, err := r.pool.Builder().
		Select(
			"uuid",
			"name",
			"symbol",
			"network_id",
			"network",
		).
		From("tokens").
		Where(goqu.C("symbol").Eq(symbol)).
		ScanStructContext(ctx, &token)

	if err != nil {
		return models.Token{}, fmt.Errorf("failed to get token: %w", err)
	}

	if !found {
		return models.Token{}, xerror.ErrNotFound
	}

	return token, nil
}
