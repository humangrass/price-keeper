package keeper

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/humangrass/price-keeper/pgk/x/xhttp"
)

func (uc *UseCase) handlePair(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = uc.getPair(w, r)
	case http.MethodPost:
		err = uc.activatePair(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	if err != nil {
		uc.logger.Sugar().Error(err)
	}
}

type priceResponse struct {
	Ticket  string                 `json:"ticket"`
	Network map[string][]priceData `json:"network"`
}

type priceData struct {
	Price     float64   `json:"price"`
	Time      time.Time `json:"time"`
	Timestamp int64     `json:"timestamp"`
}

func (uc *UseCase) getPair(w http.ResponseWriter, r *http.Request) error {
	symbol := strings.TrimPrefix(r.URL.Path, "/api/pairs/")

	parts := strings.Split(symbol, "/")
	if len(parts) != 2 {
		return xhttp.RespondWithError(w, http.StatusBadRequest, "Invalid symbol format")
	}
	prices, err := uc.priceRepository.FindBySymbols(context.Background(), parts[0], parts[1])
	if err != nil {
		uc.logger.Sugar().Error(err)
		err = xhttp.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tokens")
		return err
	}

	response := priceResponse{
		Ticket:  symbol,
		Network: make(map[string][]priceData),
	}

	for _, price := range prices {
		if price.Pair == nil {
			continue
		}

		network := price.Pair.Numerator.Network
		data := priceData{
			Price:     price.Price,
			Time:      price.TS,
			Timestamp: price.TS.Unix(),
		}

		response.Network[network] = append(response.Network[network], data)
	}

	return xhttp.RespondWithJSON(w, http.StatusOK, response)
}

func (uc *UseCase) activatePair(w http.ResponseWriter, r *http.Request) error {
	// TODO:
	return nil
}
