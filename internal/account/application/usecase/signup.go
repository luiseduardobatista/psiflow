package usecase

import (
	"github.com/google/uuid"
	"github.com/luiseduardobatista/psiflow/internal/account/domain"
)

type Signup struct{}

func (s *Signup) Execute(input SignupInput) (*SignupOuput, error) {
	account, err := domain.NewAccount(input.Email, input.Name, input.Phone, input.Password)
	if err != nil {
		return nil, err
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
