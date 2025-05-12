package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/humangrass/price-keeper/domain/entities"
	"github.com/humangrass/price-keeper/domain/models"
	"github.com/humangrass/price-keeper/pgk/xhttp"
)

func (uc *UseCase) handleTokens(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = uc.getTokens(w, r)
	case http.MethodPost:
		err = uc.createToken(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	if err != nil {
		uc.logger.Sugar().Error(err)
	}
}

func (uc *UseCase) getTokens(w http.ResponseWriter, r *http.Request) error {
	var err error
	params := entities.RequestParams{}
	params, err = params.Parse(r)
	if err != nil {
		uc.logger.Sugar().Error(err)
		err = xhttp.RespondWithError(w, http.StatusBadRequest, "Invalid request parameters")
		return err
	}

	tokens, total, err := uc.tokenRepository.GetTokens(context.Background(), params)
	if err != nil {
		uc.logger.Sugar().Error(err)
		err = xhttp.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tokens")
		return err
	}

	response := NewTokensResponse(total, tokens, params)
	return xhttp.RespondWithJSON(w, http.StatusOK, response)
}

func (uc *UseCase) createToken(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return xhttp.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

	var req NewTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return xhttp.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf(
				"Field %s: %s", err.Field(), err.Tag()))
		}
		return xhttp.RespondWithError(w, http.StatusBadRequest,
			"Validation error: "+strings.Join(validationErrors, ", "))
	}

	token := models.Token{
		Name:      req.Name,
		Symbol:    req.Symbol,
		NetworkID: req.NetworkID,
		Network:   req.Network,
	}

	if err := uc.tokenRepository.CreateToken(r.Context(), &token); err != nil {
		uc.logger.Sugar().Errorf("Failed to create token: %v", err)
		return xhttp.RespondWithError(w, http.StatusInternalServerError,
			"Failed to create token")
	}

	return xhttp.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"status": "success",
		"token":  token,
	})
}
