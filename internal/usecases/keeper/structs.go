package keeper

import (
	"fmt"
	"time"

	"github.com/humangrass/price-keeper/domain/entities"
	"github.com/humangrass/price-keeper/domain/models"
)

type TokensResponse struct {
	Data  []Token `json:"data"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
}

type Token struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	NetworkID string `json:"network_id"`
	Network   string `json:"network"`
}

func NewTokensResponse(total int, models []models.Token, params entities.RequestParams) *TokensResponse {
	page := (params.Offset / params.Limit) + 1
	response := TokensResponse{
		Total: total,
		Page:  page,
	}
	response.FillData(models)
	return &response
}

func (r *TokensResponse) FillData(models []models.Token) {
	r.Data = make([]Token, len(models))
	for i, model := range models {
		r.Data[i] = Token{
			Name:      model.Name,
			Symbol:    model.Symbol,
			NetworkID: model.NetworkID,
			Network:   model.Network,
		}
	}
}

type PairsResponse struct {
	Data  []Pair `json:"data"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
}

type Pair struct {
	Ticket    string `json:"ticket"`
	Timeframe string `json:"timeframe"`
}

func NewPairsResponse(total int, models []models.Pair, params entities.RequestParams) *PairsResponse {
	page := (params.Offset / params.Limit) + 1
	response := PairsResponse{
		Total: total,
		Page:  page,
	}
	response.FillData(models)
	return &response
}

func (r *PairsResponse) FillData(models []models.Pair) {
	r.Data = make([]Pair, len(models))
	for i, model := range models {
		r.Data[i] = Pair{
			Ticket:    fmt.Sprintf("%s/%s", model.Numerator, model.Denominator),
			Timeframe: model.Timeframe.String(),
		}
	}
}

type NewTokenRequest struct {
	Name      string `json:"name" validate:"required,max=100"`
	Symbol    string `json:"symbol" validate:"required,max=10"`
	NetworkID string `json:"network_id" validate:"required,max=100"`
	Network   string `json:"network" validate:"required,max=100"`
}

type NewPairRequest struct {
	Numerator   string        `json:"numerator" validate:"required,max=10"`
	Denominator string        `json:"denominator" validate:"required,max=10"`
	Timeframe   time.Duration `json:"timeframe" validate:"required,max=5"`
}
