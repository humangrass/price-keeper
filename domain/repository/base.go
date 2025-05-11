package repository

import "github.com/humangrass/gommon/database"

type BaseRepository struct {
	pool database.Pool
}

func NewBaseRepository(pool database.Pool) BaseRepository {
	return BaseRepository{pool: pool}
}
