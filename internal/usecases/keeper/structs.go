package keeper

import (
	"github.com/humangrass/price-keeper/domain/entities"
	"github.com/humangrass/price-keeper/domain/models"
)

type TokensResponse struct {
	Data  []Token `json:"data"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
}

type Token struct {
	UUID      string `json:"uuid"`
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
			UUID:      model.UUID.String(),
			Name:      model.Name,
			Symbol:    model.Symbol,
			NetworkID: model.NetworkID,
			Network:   model.Network,
		}
	}
}
