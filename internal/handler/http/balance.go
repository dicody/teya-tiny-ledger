package http

import (
	"net/http"
	"tiny-ledger/internal/app"
)

type balanceResponse struct {
	AccountID string `json:"account_id"`
	Balance   int64  `json:"balance"`
}

func getBalance(svc app.Ledger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := r.PathValue("id")
		if accountID == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "account id is required"})
			return
		}

		balance, err := svc.GetBalance(r.Context(), accountID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, balanceResponse{
			AccountID: accountID,
			Balance:   balance,
		})
	}
}
