package keeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// @Summary Get tokens list
// @Description Get paginated list of tokens
// @Tags tokens
// @Accept  json
// @Produce  json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Param orderBy query string false "Order by field" Enums(asc, desc) default(asc)
// @Success 200 {object} []TokensResponse
// @Failure 400 {object} xhttp.ErrorResponse
// @Failure 500 {object} xhttp.ErrorResponse
// @Router /tokens [get]
func (uc *UseCase) getTokens(w http.ResponseWriter, r *http.Request) error {
	var err error
	params := entities.RequestParams{}
	params, err = params.Parse(r)
	if err != nil {
		uc.logger.Sugar().Error(err)
		err = xhttp.RespondWithError(w, http.StatusBadRequest, "Invalid request parameters")
		return err
	}

	tokens, total, err := uc.tokenRepository.GetByParams(context.Background(), params)
	if err != nil {
		uc.logger.Sugar().Error(err)
		err = xhttp.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tokens")
		return err
	}

	response := NewTokensResponse(total, tokens, params)
	return xhttp.RespondWithJSON(w, http.StatusOK, response)
}

// @Summary Create new token
// @Description Add a new token to the database
// @Tags tokens
// @Accept  json
// @Produce  json
// @Param token body NewTokenRequest true "Token data"
// @Success 201 {object} xhttp.ErrorResponse
// @Failure 400 {object} xhttp.ErrorResponse
// @Failure 500 {object} xhttp.ErrorResponse
// @Router /tokens [post]
func (uc *UseCase) createToken(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return xhttp.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return xhttp.RespondWithError(w, http.StatusBadRequest,
			"Failed to read request body")
	}

	var req NewTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			return xhttp.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Malformed JSON at position %d", syntaxErr.Offset))
		case errors.As(err, &unmarshalErr):
			return xhttp.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Invalid value type for field '%s'", unmarshalErr.Field))
		default:
			return xhttp.RespondWithError(w, http.StatusBadRequest,
				"Invalid request body")
		}
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

	if err := uc.tokenRepository.Create(r.Context(), &token); err != nil {
		uc.logger.Sugar().Errorf("Failed to create token: %v", err)
		return xhttp.RespondWithError(w, http.StatusInternalServerError,
			"Failed to create token")
	}

	return xhttp.RespondWithJSON(w, http.StatusCreated, xhttp.ErrorResponse{
		Error: "",
	})
}
