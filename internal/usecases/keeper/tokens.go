package keeper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/humangrass/price-keeper/domain/entities"
)

func (uc *UseCase) handleTokens(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uc.getTokens(w, r)
	case http.MethodPost:
		uc.createToken(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (uc *UseCase) getTokens(w http.ResponseWriter, r *http.Request) {
	var err error
	params := entities.RequestParams{}
	params, err = params.Parse(r)
	if err != nil {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		uc.logger.Sugar().Errorf("Failed to parse request params: %v", err)
		return
	}

	tokens, total, err := uc.tokenRepository.GetTokens(context.Background(), params)
	if err != nil {
		http.Error(w, "Failed to retrieve tokens", http.StatusInternalServerError)
		uc.logger.Sugar().Errorf("Failed to get tokens: %v", err)
		return
	}

	response := NewTokensResponse(total, tokens, params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		uc.logger.Sugar().Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (uc *UseCase) createToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		uc.logger.Sugar().Errorf("HTTP server error: %v", err)
	}
}
