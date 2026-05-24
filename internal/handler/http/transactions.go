package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"tiny-ledger/internal/app"
	"tiny-ledger/internal/domain"
)

type (
	createTransactionRequest struct {
		Type   domain.TransactionType `json:"type"`
		Amount int64                  `json:"amount"`
	}
	transaction struct {
		ID        string                 `json:"id"`
		AccountID string                 `json:"account_id"`
		Type      domain.TransactionType `json:"type"`
		Amount    int64                  `json:"amount"` // positive, in cents
		CreatedAt time.Time              `json:"created_at"`
	}
	transactionsResponse struct {
		AccountID    string        `json:"account_id"`
		Transactions []transaction `json:"transactions"`
	}
)

func createTransaction(svc app.Ledger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := r.PathValue("id")
		if accountID == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "account id is required"})
			return
		}

		var req createTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
			return
		}

		tx, err := svc.RecordTransaction(r.Context(), accountID, req.Type, req.Amount)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, domain.ErrInvalidAmount) ||
				errors.Is(err, domain.ErrInsufficientFunds) ||
				errors.Is(err, domain.ErrInvalidType) {
				status = http.StatusBadRequest
			}
			writeJSON(w, status, errorResponse{Error: err.Error()})
			return
		}

		writeJSON(w, http.StatusCreated, transaction{
			ID:        tx.ID,
			AccountID: tx.AccountID,
			Type:      tx.Type,
			Amount:    tx.Amount,
			CreatedAt: tx.CreatedAt,
		})
	}
}

func getTransactions(svc app.Ledger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := r.PathValue("id")
		if accountID == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "account id is required"})
			return
		}

		transactions, err := svc.GetTransactions(r.Context(), accountID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}

		txs := make([]transaction, len(transactions))
		for i, t := range transactions {
			txs[i] = transaction{
				ID:        t.ID,
				AccountID: t.AccountID,
				Type:      t.Type,
				Amount:    t.Amount,
				CreatedAt: t.CreatedAt,
			}
		}

		writeJSON(w, http.StatusOK, transactionsResponse{
			AccountID:    accountID,
			Transactions: txs,
		})
	}
}
