package domain

import (
	"github.com/google/uuid"
)

type AccountRepository interface {
	Save(account Account) error
	GetByID(accountID uuid.UUID) (*Account, error)
}
