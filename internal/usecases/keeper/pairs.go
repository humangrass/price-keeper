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
	"github.com/humangrass/price-keeper/pgk/x/xerror"
	"github.com/humangrass/price-keeper/pgk/x/xhttp"
	"github.com/humangrass/price-keeper/pgk/x/xtype"
)

func (uc *UseCase) handlePairs(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = uc.getPairs(w, r)
	case http.MethodPost:
		err = uc.createPair(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	if err != nil {
		uc.logger.Sugar().Error(err)
	}
}

// getPairs retrieves token pairs
// @Summary List token pairs
// @Description Get paginated list of token pairs with sorting options
// @Tags pairs
// @Produce json
// @Param offset query int false "Pagination offset" minimum(0) default(0)
// @Param limit query int false "Items per page" minimum(1) maximum(100) default(10)
// @Param orderBy query string false "Sort order" Enums(asc,desc) default(asc)
// @Success 200 {object} PairsResponse "Successful operation"
// @Failure 400 {object} xhttp.ErrorResponse "Invalid parameters"
// @Failure 500 {object} xhttp.ErrorResponse "Internal server error"
// @Router /pairs [get]
func (uc *UseCase) getPairs(w http.ResponseWriter, r *http.Request) error {
	var err error
	params := entities.RequestParams{}
	params, err = params.Parse(r)
	if err != nil {
		uc.logger.Sugar().Error(err)
		err = xhttp.RespondWithError(w, http.StatusBadRequest, "Invalid request parameters")
		return err
	}

	pairs, total, err := uc.pairsRepository.GetByParams(context.Background(), params)
	if err != nil {
		uc.logger.Sugar().Error(err)
		err = xhttp.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tokens")
		return err
	}

	response := NewPairsResponse(total, pairs, params)
	return xhttp.RespondWithJSON(w, http.StatusOK, response)
}

// createPair creates new token pair
// @Summary Create token pair
// @Description Register new token pair in the system
// @Tags pairs
// @Accept json
// @Produce json
// @Param pair body NewPairRequest true "Token pair data"
// @Success 201 {object} Pair "Successfully created"
// @Failure 400 {object} xhttp.ErrorResponse "Invalid input"
// @Failure 404 {object} xhttp.ErrorResponse "Token not found"
// @Failure 500 {object} xhttp.ErrorResponse "Internal server error"
// @Router /pairs [post]
func (uc *UseCase) createPair(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return xhttp.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return xhttp.RespondWithError(w, http.StatusBadRequest,
			"Failed to read request body")
	}

	var req NewPairRequest
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
	if err = validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf(
				"Field %s: %s", err.Field(), err.Tag()))
		}
		return xhttp.RespondWithError(w, http.StatusBadRequest,
			"Validation error: "+strings.Join(validationErrors, ", "))
	}

	numerator, err := uc.tokenRepository.GetTokenBySymbol(r.Context(), req.Numerator)
	if err != nil {
		if xerror.IsNotFound(err) {
			return xhttp.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Numerator token '%s' not found", req.Numerator))
		}
		uc.logger.Sugar().Errorf("Failed to get numerator token: %v", err)
		return xhttp.RespondWithError(w, http.StatusInternalServerError,
			"Failed to create pair")
	}
	denominator, err := uc.tokenRepository.GetTokenBySymbol(r.Context(), req.Denominator)
	if err != nil {
		if xerror.IsNotFound(err) {
			return xhttp.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Denominator token '%s' not found", req.Denominator))
		}
		uc.logger.Sugar().Errorf("Failed to get denominator token: %v", err)
		return xhttp.RespondWithError(w, http.StatusInternalServerError,
			"Failed to create pair")
	}

	pair := models.Pair{
		Numerator:   models.Token{UUID: numerator.UUID},
		Denominator: models.Token{UUID: denominator.UUID},
		Timeframe:   xtype.FromDuration(req.Timeframe),
		IsActive:    false,
	}

	if err := uc.pairsRepository.Create(r.Context(), &pair); err != nil {
		uc.logger.Sugar().Errorf("Failed to create pair: %v", err)
		return xhttp.RespondWithError(w, http.StatusInternalServerError,
			"Failed to create pair")
	}

	return xhttp.RespondWithJSON(w, http.StatusCreated, xhttp.ErrorResponse{
		Error: "",
	})
}
