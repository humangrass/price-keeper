package keeper

import (
	"net/http"
	"strings"

	"github.com/humangrass/price-keeper/pgk/xhttp"
)

func (uc *UseCase) handlePair(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = uc.getPair(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	if err != nil {
		uc.logger.Sugar().Error(err)
	}
}

func (uc *UseCase) getPair(w http.ResponseWriter, r *http.Request) error {
	symbol := strings.TrimPrefix(r.URL.Path, "/api/pairs/")

	parts := strings.Split(symbol, "/")
	if len(parts) != 2 {
		return xhttp.RespondWithError(w, http.StatusBadRequest, "Invalid symbol format")
	}

	return nil
}
