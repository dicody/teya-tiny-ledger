package http

import (
	"encoding/json"
	"net/http"
	"tiny-ledger/internal/app"
)

type errorResponse struct {
	Error string `json:"error"`
}

func NewHTTPHandler(svc app.Ledger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/accounts/{id}/transactions", getTransactions(svc))
	mux.HandleFunc("POST /api/v1/accounts/{id}/transactions", createTransaction(svc))
	mux.HandleFunc("GET /api/v1/accounts/{id}/balance", getBalance(svc))

	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
