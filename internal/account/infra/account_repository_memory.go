package infra

import (
	"sync"

	"github.com/google/uuid"
	"github.com/luiseduardobatista/psiflow/internal/account/domain"
)

type AccountRepositoryMemory struct {
	mu       sync.RWMutex
	accounts map[string]domain.Account
}

var _ domain.AccountRepository = (*AccountRepositoryMemory)(nil)

func NewAccountRepositoryMemory() *AccountRepositoryMemory {
	return &AccountRepositoryMemory{
		accounts: make(map[string]domain.Account),
	}
}

func (a *AccountRepositoryMemory) Save(account domain.Account) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.accounts[account.ID.String()] = account
	return nil
}

func (a *AccountRepositoryMemory) GetByID(accountID uuid.UUID) (*domain.Account, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	account, exists := a.accounts[accountID.String()]
	if !exists {
		return nil, nil
	}
	return &account, nil
}
