package keeper

import (
	"github.com/humangrass/price-keeper/pgk/logger"
)

type UseCase struct {
	logger *logger.Logger
}

func NewKeeperUseCase(
	logger *logger.Logger,
) *UseCase {
	uc := &UseCase{
		logger: logger,
	}

	return uc
}
