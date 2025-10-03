package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrEmailAlreadyExists = errors.New("account with this email already exists")
)

type AccountRepository interface {
	Save(account *Account) error
	GetByID(accountID uuid.UUID) (*Account, error)
	Update(account *Account) error
}
