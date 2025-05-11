package jupiter

//easyjson:json
type PriceResponse struct {
	Data      map[string]PriceData `json:"data"`
	TimeTaken float64              `json:"timeTaken"`
}

type PriceData struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Price string `json:"price"`
}
