package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tiny-ledger/internal/app"
	"tiny-ledger/internal/domain"
	"tiny-ledger/internal/storage"
)

func setupTestServer() *httptest.Server {
	store := storage.NewMemStore()
	svc := app.NewLedger(store)
	router := NewHTTPHandler(svc)
	return httptest.NewServer(router)
}

func TestCreateDeposit(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	body := `{"type":"deposit","amount":5000}`
	resp, err := http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var tx transaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if tx.AccountID != "dmitry" || tx.Type != domain.Deposit || tx.Amount != 5000 {
		t.Errorf("unexpected transaction: %+v", tx)
	}
}

func TestCreateWithdrawal(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	deposit := `{"type":"deposit","amount":5000}`
	_, _ = http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json", bytes.NewBufferString(deposit))

	withdraw := `{"type":"withdrawal","amount":2000}`
	resp, err := http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json", bytes.NewBufferString(withdraw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var tx transaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if tx.Type != domain.Withdrawal || tx.Amount != 2000 {
		t.Errorf("unexpected transaction: %+v", tx)
	}
}

func TestWithdrawalInsufficientFunds(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	body := `{"type":"withdrawal","amount":1000}`
	resp, err := http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestInvalidAmount(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	body := `{"type":"deposit","amount":0}`
	resp, err := http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestInvalidTransactionType(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	body := `{"type":"transfer","amount":1000}`
	resp, err := http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetBalance(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	deposit := `{"type":"deposit","amount":5000}`
	_, _ = http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json", bytes.NewBufferString(deposit))

	withdraw := `{"type":"withdrawal","amount":1500}`
	_, _ = http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json", bytes.NewBufferString(withdraw))

	resp, err := http.Get(srv.URL + "/api/v1/accounts/dmitry/balance")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result balanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.AccountID != "dmitry" || result.Balance != 3500 {
		t.Errorf("unexpected balance response: %+v", result)
	}
}

func TestGetBalanceEmptyAccount(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/accounts/unknown/balance")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result balanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Balance != 0 {
		t.Errorf("expected balance 0, got %d", result.Balance)
	}
}

func TestGetTransactions(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	deposit := `{"type":"deposit","amount":3000}`
	_, _ = http.Post(srv.URL+"/api/v1/accounts/bob/transactions", "application/json", bytes.NewBufferString(deposit))

	withdraw := `{"type":"withdrawal","amount":1000}`
	_, _ = http.Post(srv.URL+"/api/v1/accounts/bob/transactions", "application/json", bytes.NewBufferString(withdraw))

	resp, err := http.Get(srv.URL + "/api/v1/accounts/bob/transactions")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result transactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.AccountID != "bob" {
		t.Errorf("expected account_id 'bob', got %q", result.AccountID)
	}

	if len(result.Transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(result.Transactions))
	}

	if result.Transactions[0].Type != domain.Deposit || result.Transactions[0].Amount != 3000 {
		t.Errorf("unexpected first transaction: %+v", result.Transactions[0])
	}

	if result.Transactions[1].Type != domain.Withdrawal || result.Transactions[1].Amount != 1000 {
		t.Errorf("unexpected second transaction: %+v", result.Transactions[1])
	}
}

func TestInvalidRequestBody(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/v1/accounts/dmitry/transactions", "application/json",
		bytes.NewBufferString(`not json`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}
