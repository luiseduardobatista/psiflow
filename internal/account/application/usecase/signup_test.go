package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/luiseduardobatista/psiflow/internal/account/application/usecase"
	"github.com/luiseduardobatista/psiflow/internal/account/infra"
)

func TestShouldCreateValidAccount(t *testing.T) {
	assert := require.New(t)
	inputSignup := usecase.SignupInput{
		Name:     "John Doe",
		Email:    fmt.Sprintf("user%d@example.com", time.Now().UnixNano()),
		Password: "asdQWE123",
		Phone:    "+5511999999999",
	}
	accountRepository := infra.NewAccountRepositoryMemory()
	signup := usecase.NewSignupUseCase(accountRepository)
	signupOutput, err := signup.Execute(inputSignup)
	assert.NoError(err)
	assert.NotNil(signupOutput)
	assert.NotEmptyf(signupOutput.AccountID, "AccountID should not be empty")
	getAccount := usecase.NewGetAccountUseCase(accountRepository)
	getAccountOutput, err := getAccount.Execute(signupOutput.AccountID)
	assert.NoError(err)
	assert.Equal(getAccountOutput.Name.String(), inputSignup.Name)
	assert.Equal(getAccountOutput.Email.String(), inputSignup.Email)
	assert.Equal(getAccountOutput.Phone.String(), inputSignup.Phone)
}
