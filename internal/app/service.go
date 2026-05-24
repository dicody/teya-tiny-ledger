package app

import (
	"context"
	"tiny-ledger/internal/domain"
)

type Repository interface {
	AppendTransaction(ctx context.Context, tx domain.Transaction) error
	ListTransactions(ctx context.Context, accountID string) ([]domain.Transaction, error)
}

type Ledger interface {
	RecordTransaction(ctx context.Context, accountID string, txType domain.TransactionType, amount int64) (domain.Transaction, error)
	GetBalance(ctx context.Context, accountID string) (int64, error)
	GetTransactions(ctx context.Context, accountID string) ([]domain.Transaction, error)
}
