package repository

import "github.com/humangrass/gommon/database"

type TokensRepository struct {
	BaseRepository
}

func NewTokensRepository(pool database.Pool) *TokensRepository {
	return &TokensRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

type TokenRepo interface {
	// TODO:
}
