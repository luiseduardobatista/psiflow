package usecase

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/luiseduardobatista/psiflow/internal/account/domain"
	"github.com/luiseduardobatista/psiflow/internal/shared"
)

type Signup struct {
	accountRepository domain.AccountRepository
}

func NewSignupUseCase(accountRepository domain.AccountRepository) *Signup {
	return &Signup{
		accountRepository: accountRepository,
	}
}

func (s *Signup) Execute(input SignupInput) (*SignupOuput, error) {
	account, err := domain.NewAccount(input.Email, input.Name, input.Phone, input.Password)
	if err != nil {
		return nil, shared.NewDomainError(http.StatusUnprocessableEntity, err.Error())
	}
	err = s.accountRepository.Save(account)
	if err != nil {
		return nil, shared.NewInfraError(err)
	}
	return &SignupOuput{
		AccountID: account.ID,
	}, nil
}

type SignupInput struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

type SignupOuput struct {
	AccountID uuid.UUID
}
