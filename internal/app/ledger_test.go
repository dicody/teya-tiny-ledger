package app

import (
	"context"
	"errors"
	"testing"
	"tiny-ledger/internal/domain"
	"tiny-ledger/internal/storage"
)

func TestRecordTransaction_Deposit(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())

	tx, err := svc.RecordTransaction(context.Background(), "dmitry", domain.Deposit, 5000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tx.AccountID != "dmitry" || tx.Type != domain.Deposit || tx.Amount != 5000 {
		t.Errorf("unexpected transaction: %+v", tx)
	}

	if tx.ID == "" {
		t.Error("expected non-empty transaction ID")
	}
}

func TestRecordTransaction_Withdrawal(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())
	ctx := context.Background()

	_, _ = svc.RecordTransaction(ctx, "foo", domain.Deposit, 5000)

	tx, err := svc.RecordTransaction(ctx, "foo", domain.Withdrawal, 2000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tx.AccountID != "foo" || tx.Type != domain.Withdrawal || tx.Amount != 2000 {
		t.Errorf("unexpected transaction: %+v", tx)
	}

	balance, err := svc.GetBalance(ctx, "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance != 3000 {
		t.Errorf("expected balance 3000, got %d", balance)
	}
}

func TestRecordTransaction_InvalidType(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())

	_, err := svc.RecordTransaction(context.Background(), "dmitry", "transfer", 1000)
	if !errors.Is(err, domain.ErrInvalidType) {
		t.Errorf("expected ErrInvalidType, got %v", err)
	}
}

func TestRecordTransaction_InvalidAmount(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())

	_, err := svc.RecordTransaction(context.Background(), "dmitry", domain.Deposit, 0)
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}

	_, err = svc.RecordTransaction(context.Background(), "dmitry", domain.Withdrawal, -500)
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestRecordTransaction_InsufficientFunds(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())
	ctx := context.Background()

	_, err := svc.RecordTransaction(ctx, "foo", domain.Withdrawal, 1000)
	if !errors.Is(err, domain.ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}

func TestRecordTransaction_ExactBalance(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())
	ctx := context.Background()

	_, _ = svc.RecordTransaction(ctx, "foo", domain.Deposit, 3000)

	_, err := svc.RecordTransaction(ctx, "foo", domain.Withdrawal, 3000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	balance, err := svc.GetBalance(ctx, "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance != 0 {
		t.Errorf("expected balance 0, got %d", balance)
	}
}

func TestGetBalance_UpdatesAfterDeposits(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())
	ctx := context.Background()

	_, _ = svc.RecordTransaction(ctx, "dmitry", domain.Deposit, 1000)
	_, _ = svc.RecordTransaction(ctx, "dmitry", domain.Deposit, 2500)

	balance, err := svc.GetBalance(ctx, "dmitry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if balance != 3500 {
		t.Errorf("expected balance 3500, got %d", balance)
	}
}

func TestGetBalance_NonExistentAccount(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())

	balance, err := svc.GetBalance(context.Background(), "unknown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if balance != 0 {
		t.Errorf("expected balance 0 for non-existent account, got %d", balance)
	}
}

func TestGetTransactions(t *testing.T) {
	svc := NewLedger(storage.NewMemStore())
	ctx := context.Background()

	_, err := svc.RecordTransaction(ctx, "dmitry", domain.Deposit, 3000)
	if err != nil {
		t.Fatalf("deposit: %v", err)
	}
	_, err = svc.RecordTransaction(ctx, "dmitry", domain.Withdrawal, 1000)
	if err != nil {
		t.Fatalf("withdraw: %v", err)
	}

	txs, err := svc.GetTransactions(ctx, "dmitry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(txs) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(txs))
	}

	if txs[0].Type != domain.Deposit || txs[0].Amount != 3000 {
		t.Errorf("unexpected first transaction: %+v", txs[0])
	}
	if txs[1].Type != domain.Withdrawal || txs[1].Amount != 1000 {
		t.Errorf("unexpected second transaction: %+v", txs[1])
	}
}
