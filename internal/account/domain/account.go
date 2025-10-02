package domain

import (
	"github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID
	Name     Name
	Email    Email
	Phone    Phone
	password Password
}

func NewAccount(email string, name string, phone string, password string) (*Account, error) {
	accountID := uuid.New()
	accountEmail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	accountName, err := NewName(name)
	if err != nil {
		return nil, err
	}
	accountPhone, err := NewPhone(phone)
	if err != nil {
		return nil, err
	}
	accountPassword, err := NewPassword(password)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:       accountID,
		Email:    accountEmail,
		Name:     accountName,
		Phone:    accountPhone,
		password: accountPassword,
	}, nil
}
