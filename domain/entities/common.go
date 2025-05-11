package entities

import (
	"fmt"
	"net/http"
	"strconv"
)

type RequestParams struct {
	Offset  int
	Limit   int
	OrderBy OrderBy
}

func (p *RequestParams) Parse(r *http.Request) (RequestParams, error) {
	query := r.URL.Query()

	params := RequestParams{
		Offset:  0,
		Limit:   10,
		OrderBy: "asc",
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return params, fmt.Errorf("invalid offset value")
		}
		params.Offset = offset
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			return params, fmt.Errorf("limit must be between 1 and 100")
		}
		params.Limit = limit
	}

	if orderBy := query.Get("orderBy"); orderBy != "" {
		order := OrderBy(orderBy)
		if !order.isValid() {
			return params, fmt.Errorf("orderBy must be '%s' or '%s'", OrderByAsc, OrderByDesc)
		}
		params.OrderBy = order
	}

	return params, nil

}

type OrderBy string

const (
	OrderByAsc  OrderBy = "asc"
	OrderByDesc OrderBy = "desc"
)

func (o OrderBy) isValid() bool {
	return o == OrderByAsc || o == OrderByDesc
}
