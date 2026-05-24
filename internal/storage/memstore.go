package storage

import (
	"context"
	"sync"

	"tiny-ledger/internal/domain"
)

type MemStore struct {
	mu           sync.RWMutex
	transactions map[string][]domain.Transaction // accountID → []Transaction
}

func NewMemStore() *MemStore {
	return &MemStore{
		transactions: make(map[string][]domain.Transaction),
	}
}

func (m *MemStore) AppendTransaction(_ context.Context, tx domain.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transactions[tx.AccountID] = append(m.transactions[tx.AccountID], tx)
	return nil
}

func (m *MemStore) ListTransactions(_ context.Context, accountID string) ([]domain.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	txs := m.transactions[accountID]
	result := make([]domain.Transaction, len(txs))
	copy(result, txs)
	return result, nil
}
