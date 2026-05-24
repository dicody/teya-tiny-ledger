package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
	"tiny-ledger/internal/domain"
)

type ledger struct {
	repo Repository
}

func NewLedger(repo Repository) Ledger {
	return &ledger{repo: repo}
}

func (l *ledger) RecordTransaction(ctx context.Context, accountID string, txType domain.TransactionType, amount int64) (domain.Transaction, error) {
	switch txType {
	case domain.Deposit:
		return l.deposit(ctx, accountID, amount)
	case domain.Withdrawal:
		return l.withdraw(ctx, accountID, amount)
	default:
		return domain.Transaction{}, domain.ErrInvalidType
	}
}

func (l *ledger) GetBalance(ctx context.Context, accountID string) (int64, error) {
	transactions, err := l.repo.ListTransactions(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("list transactions: %w", err)
	}

	balance, err := domain.CalculateBalance(transactions)
	if err != nil {
		return 0, fmt.Errorf("calculate balance: %w", err)
	}

	return balance, nil
}

func (l *ledger) deposit(ctx context.Context, accountID string, amount int64) (domain.Transaction, error) {
	tx, err := domain.NewTransaction(generateID(), accountID, domain.Deposit, amount, time.Now())
	if err != nil {
		return domain.Transaction{}, err
	}

	if err := l.repo.AppendTransaction(ctx, tx); err != nil {
		return domain.Transaction{}, fmt.Errorf("append transaction: %w", err)
	}

	return tx, nil
}

// withdraw records a withdrawal. Note: there is a TOCTOU race between
// balance check and append — concurrent withdrawals could overdraw.
// See README for rationale on why this is accepted.
func (l *ledger) withdraw(ctx context.Context, accountID string, amount int64) (domain.Transaction, error) {
	tx, err := domain.NewTransaction(generateID(), accountID, domain.Withdrawal, amount, time.Now())
	if err != nil {
		return domain.Transaction{}, err
	}

	transactions, err := l.repo.ListTransactions(ctx, accountID)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("list transactions: %w", err)
	}

	balance, err := domain.CalculateBalance(transactions)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("calculate balance: %w", err)
	}

	if balance < amount {
		return domain.Transaction{}, domain.ErrInsufficientFunds
	}

	if err := l.repo.AppendTransaction(ctx, tx); err != nil {
		return domain.Transaction{}, fmt.Errorf("append transaction: %w", err)
	}

	return tx, nil
}

func (l *ledger) GetTransactions(ctx context.Context, accountID string) ([]domain.Transaction, error) {
	transactions, err := l.repo.ListTransactions(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}

	return transactions, nil
}

func generateID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand: " + err.Error())
	}
	return fmt.Sprintf("%x", b)
}
