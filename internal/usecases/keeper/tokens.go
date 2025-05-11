package keeper

import "net/http"

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
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("TOKENS!"))
	if err != nil {
		uc.logger.Sugar().Errorf("HTTP server error: %v", err)
	}

}

func (uc *UseCase) createToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		uc.logger.Sugar().Errorf("HTTP server error: %v", err)
	}
}
