package repository

import (
	"github.com/humangrass/gommon/database"
)

type PairsRepository struct {
	BaseRepository
}

func NewPairsRepository(pool database.Pool) *PairsRepository {
	return &PairsRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

type PairsRepo interface {
	// TODO:
	// Create(ctx context.Context) (models.Pair, error)
}
