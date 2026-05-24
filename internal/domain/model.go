package domain

import (
	"errors"
	"fmt"
	"time"
)

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)

type Transaction struct {
	ID        string
	AccountID string
	Type      TransactionType
	Amount    int64 // positive, in cents
	CreatedAt time.Time
}

var (
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInvalidType       = errors.New("transaction type must be 'deposit' or 'withdrawal'")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

func NewTransaction(id, accountID string, txType TransactionType, amount int64, createdAt time.Time) (Transaction, error) {
	if err := ValidateAmount(amount); err != nil {
		return Transaction{}, err
	}
	if err := ValidateType(txType); err != nil {
		return Transaction{}, err
	}

	return Transaction{
		ID:        id,
		AccountID: accountID,
		Type:      txType,
		Amount:    amount,
		CreatedAt: createdAt,
	}, nil
}

func CalculateBalance(transactions []Transaction) (int64, error) {
	var balance int64
	for _, tx := range transactions {
		switch tx.Type {
		case Deposit:
			balance += tx.Amount
		case Withdrawal:
			balance -= tx.Amount
		default:
			return 0, fmt.Errorf("unknown transaction type: %q", tx.Type)
		}
	}
	return balance, nil
}

func ValidateAmount(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

func ValidateType(t TransactionType) error {
	switch t {
	case Deposit, Withdrawal:
		return nil
	default:
		return ErrInvalidType
	}
}
