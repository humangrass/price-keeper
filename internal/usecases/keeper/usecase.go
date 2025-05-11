package keeper

import (
	"net/http"

	"github.com/humangrass/price-keeper/domain/repository"
	"github.com/humangrass/price-keeper/pgk/logger"
)

type UseCase struct {
	logger *logger.Logger

	pairsRepository repository.PairsRepo
	priceRepository repository.PriceRepo
	tokenRepository repository.TokenRepo
}

func (uc *UseCase) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/tokens", uc.handleTokens)
}

func NewKeeperUseCase(
	baseRepo *repository.BaseRepository,
	tokenRepo *repository.TokensRepository,
	logger *logger.Logger,
) *UseCase {
	uc := &UseCase{
		logger:          logger,
		pairsRepository: baseRepo,
		priceRepository: baseRepo,
		tokenRepository: tokenRepo,
	}

	return uc
}
