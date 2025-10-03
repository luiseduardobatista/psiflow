package usecase

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/luiseduardobatista/psiflow/internal/account/domain"
	"github.com/luiseduardobatista/psiflow/internal/shared"
)

type UpdateAccount struct {
	accountRepository domain.AccountRepository
}

func NewUpdateAccountUseCase(accountRepository domain.AccountRepository) *UpdateAccount {
	return &UpdateAccount{
		accountRepository: accountRepository,
	}
}

func (u *UpdateAccount) Execute(input UpdateAccountInput) (*UpdateAccountOuput, error) {
	account, err := u.accountRepository.GetByID(input.ID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, shared.NewDomainError(http.StatusNotFound, err.Error())
		}
		return nil, shared.NewInfraError(err)
	}
	if err := u.applyChanges(account, input); err != nil {
		return nil, err
	}
	if err := u.accountRepository.Update(account); err != nil {
		return nil, shared.NewInfraError(err)
	}
	return &UpdateAccountOuput{
		ID:    account.ID,
		Name:  account.Name.String(),
		Email: account.Email.String(),
		Phone: account.Phone.String(),
	}, nil
}

func (u *UpdateAccount) applyChanges(account *domain.Account, input UpdateAccountInput) error {
	if input.Name != "" {
		if err := account.ChangeName(input.Name); err != nil {
			return shared.NewDomainError(http.StatusUnprocessableEntity, err.Error())
		}
	}
	if input.Phone != "" {
		if err := account.ChangePhone(input.Phone); err != nil {
			return shared.NewDomainError(http.StatusUnprocessableEntity, err.Error())
		}
	}
	return nil
}

type UpdateAccountInput struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
	Phone    string
}

type UpdateAccountOuput struct {
	ID    uuid.UUID
	Name  string
	Email string
	Phone string
}
