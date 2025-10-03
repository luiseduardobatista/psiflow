package usecase

import (
	"github.com/google/uuid"
	"github.com/luiseduardobatista/psiflow/internal/account/domain"
	"github.com/luiseduardobatista/psiflow/internal/shared"
)

type GetAccount struct {
	accountRepository domain.AccountRepository
}

func NewGetAccountUseCase(accountRepository domain.AccountRepository) *GetAccount {
	return &GetAccount{
		accountRepository: accountRepository,
	}
}

func (s *GetAccount) Execute(input uuid.UUID) (*GetAccountOuput, error) {
	account, err := s.accountRepository.GetByID(input)
	if err != nil {
		return nil, shared.NewInfraError(err)
	}
	return &GetAccountOuput{
		ID:    account.ID,
		Name:  account.Name,
		Email: account.Email,
		Phone: account.Phone,
	}, nil
}

type GetAccountOuput struct {
	ID    uuid.UUID
	Name  domain.Name
	Email domain.Email
	Phone domain.Phone
}
