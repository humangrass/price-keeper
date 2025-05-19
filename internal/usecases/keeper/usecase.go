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
	mux.HandleFunc("/api/pairs", uc.handlePairs)
	mux.HandleFunc("/api/pairs/", uc.handlePair)
}

func NewKeeperUseCase(
	priceRepo *repository.PricesRepository,
	tokenRepo *repository.TokensRepository,
	pairsRepo *repository.PairsRepository,
	logger *logger.Logger,
) *UseCase {
	uc := &UseCase{
		logger:          logger,
		pairsRepository: pairsRepo,
		priceRepository: priceRepo,
		tokenRepository: tokenRepo,
	}

	return uc
}
