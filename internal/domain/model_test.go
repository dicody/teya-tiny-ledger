package domain

import (
	"errors"
	"testing"
	"time"
)

func TestCalculateBalance_Empty(t *testing.T) {
	balance, err := CalculateBalance(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance != 0 {
		t.Errorf("expected 0, got %d", balance)
	}
}

func TestCalculateBalance_DepositsOnly(t *testing.T) {
	txs := []Transaction{
		{Type: Deposit, Amount: 1000},
		{Type: Deposit, Amount: 2500},
	}
	balance, err := CalculateBalance(txs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance != 3500 {
		t.Errorf("expected 3500, got %d", balance)
	}
}

func TestCalculateBalance_Mixed(t *testing.T) {
	txs := []Transaction{
		{Type: Deposit, Amount: 5000},
		{Type: Withdrawal, Amount: 2000},
		{Type: Deposit, Amount: 1000},
		{Type: Withdrawal, Amount: 500},
	}
	balance, err := CalculateBalance(txs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance != 3500 {
		t.Errorf("expected 3500, got %d", balance)
	}
}

func TestCalculateBalance_ToZero(t *testing.T) {
	txs := []Transaction{
		{Type: Deposit, Amount: 1000},
		{Type: Withdrawal, Amount: 1000},
	}
	balance, err := CalculateBalance(txs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance != 0 {
		t.Errorf("expected 0, got %d", balance)
	}
}

func TestCalculateBalance_UnknownType(t *testing.T) {
	txs := []Transaction{
		{Type: Deposit, Amount: 1000},
		{Type: "refund", Amount: 500},
	}
	_, err := CalculateBalance(txs)
	if err == nil {
		t.Fatal("expected error for unknown transaction type, got nil")
	}
}

func TestValidateAmount_Positive(t *testing.T) {
	if err := ValidateAmount(100); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateAmount_Zero(t *testing.T) {
	if err := ValidateAmount(0); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestValidateAmount_Negative(t *testing.T) {
	if err := ValidateAmount(-50); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestValidateType_Valid(t *testing.T) {
	for _, tt := range []TransactionType{Deposit, Withdrawal} {
		if err := ValidateType(tt); err != nil {
			t.Errorf("expected nil for %q, got %v", tt, err)
		}
	}
}

func TestValidateType_Invalid(t *testing.T) {
	if err := ValidateType("transfer"); !errors.Is(err, ErrInvalidType) {
		t.Errorf("expected ErrInvalidType, got %v", err)
	}
}

func TestNewTransaction_Valid(t *testing.T) {
	now := time.Now()
	tx, err := NewTransaction("tx-1", "dmitry", Deposit, 1050, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.ID != "tx-1" || tx.AccountID != "dmitry" || tx.Type != Deposit || tx.Amount != 1050 || tx.CreatedAt != now {
		t.Errorf("unexpected transaction fields: %+v", tx)
	}
}

func TestNewTransaction_InvalidAmount(t *testing.T) {
	_, err := NewTransaction("tx-1", "dmitry", Deposit, 0, time.Now())
	if !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestNewTransaction_InvalidType(t *testing.T) {
	_, err := NewTransaction("tx-1", "dmitry", "transfer", 1000, time.Now())
	if !errors.Is(err, ErrInvalidType) {
		t.Errorf("expected ErrInvalidType, got %v", err)
	}
}
