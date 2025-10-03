package domain

import (
	"github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID
	Name     Name
	Email    Email
	Phone    Phone
	Address  Address
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

func (a *Account) ChangeName(newName string) error {
	nameVO, err := NewName(newName)
	if err != nil {
		return err
	}
	a.Name = nameVO
	return nil
}

func (a *Account) ChangePhone(newPhone string) error {
	phoneVO, err := NewPhone(newPhone)
	if err != nil {
		return err
	}
	a.Phone = phoneVO
	return nil
}
