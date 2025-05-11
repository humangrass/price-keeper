package repository

import "github.com/humangrass/gommon/database"

type PricesRepository struct {
	BaseRepository
}

func NewPricesRepository(pool database.Pool) *PricesRepository {
	return &PricesRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

type PriceRepo interface {
	// TODO:
}
