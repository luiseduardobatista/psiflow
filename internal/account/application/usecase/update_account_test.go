package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luiseduardobatista/psiflow/internal/account/application/usecase"
	"github.com/luiseduardobatista/psiflow/internal/account/infra"
	"github.com/stretchr/testify/require"
)

func setupTestAccount(t *testing.T, repo *infra.AccountRepositoryMemory) (*usecase.SignupOuput, *usecase.SignupInput) {
	inputSignup := &usecase.SignupInput{
		Name:     "John Original",
		Email:    fmt.Sprintf("user%d@example.com", time.Now().UnixNano()),
		Password: "asdQWE123",
		Phone:    "+5511988888888",
	}
	signup := usecase.NewSignupUseCase(repo)
	signupOutput, err := signup.Execute(*inputSignup)
	require.NoError(t, err)
	return signupOutput, inputSignup
}

func TestShouldUpdateAccountNameAndPhone(t *testing.T) {
	assert := require.New(t)
	accountRepository := infra.NewAccountRepositoryMemory()
	signupOutput, originalInput := setupTestAccount(t, accountRepository)
	updateInput := usecase.UpdateAccountInput{
		ID:    signupOutput.AccountID,
		Name:  "John Updated",
		Phone: "+5511977777777",
	}
	updateAccount := usecase.NewUpdateAccountUseCase(accountRepository)
	updateOutput, err := updateAccount.Execute(updateInput)
	assert.NoError(err)
	assert.NotNil(updateOutput)
	assert.Equal(updateInput.Name, updateOutput.Name)
	assert.Equal(updateInput.Phone, updateOutput.Phone)
	assert.Equal(originalInput.Email, updateOutput.Email)
	getAccount := usecase.NewGetAccountUseCase(accountRepository)
	getAccountOutput, err := getAccount.Execute(signupOutput.AccountID)
	assert.NoError(err)
	assert.Equal(updateInput.Name, getAccountOutput.Name.String())
	assert.Equal(updateInput.Phone, getAccountOutput.Phone.String())
	assert.Equal(originalInput.Email, getAccountOutput.Email.String())
}

func TestShouldUpdateOnlyAccountName(t *testing.T) {
	assert := require.New(t)
	accountRepository := infra.NewAccountRepositoryMemory()
	signupOutput, originalInput := setupTestAccount(t, accountRepository)
	updateInput := usecase.UpdateAccountInput{
		ID:   signupOutput.AccountID,
		Name: "Jane Doe Updated",
	}
	updateAccount := usecase.NewUpdateAccountUseCase(accountRepository)
	updateOutput, err := updateAccount.Execute(updateInput)
	assert.NoError(err)
	assert.NotNil(updateOutput)
	getAccount := usecase.NewGetAccountUseCase(accountRepository)
	getAccountOutput, err := getAccount.Execute(signupOutput.AccountID)
	assert.NoError(err)
	assert.Equal(updateInput.Name, getAccountOutput.Name.String())
	assert.Equal(originalInput.Phone, getAccountOutput.Phone.String())
}

func TestShouldReturnErrorWhenUpdatingNonExistentAccount(t *testing.T) {
	assert := require.New(t)
	accountRepository := infra.NewAccountRepositoryMemory()
	updateInput := usecase.UpdateAccountInput{
		ID:   uuid.New(),
		Name: "Some Name",
	}
	updateAccount := usecase.NewUpdateAccountUseCase(accountRepository)
	updateOutput, err := updateAccount.Execute(updateInput)
	assert.Error(err)
	assert.Nil(updateOutput)
	assert.Contains(err.Error(), "account not found")
}

func TestShouldReturnErrorWhenUpdatingWithInvalidName(t *testing.T) {
	assert := require.New(t)
	accountRepository := infra.NewAccountRepositoryMemory()
	signupOutput, originalInput := setupTestAccount(t, accountRepository)
	updateInput := usecase.UpdateAccountInput{
		ID:   signupOutput.AccountID,
		Name: "Invalid",
	}
	updateAccount := usecase.NewUpdateAccountUseCase(accountRepository)
	updateOutput, err := updateAccount.Execute(updateInput)
	assert.Error(err)
	assert.Nil(updateOutput)
	assert.Contains(err.Error(), "invalid name")
	getAccount := usecase.NewGetAccountUseCase(accountRepository)
	getAccountOutput, err := getAccount.Execute(signupOutput.AccountID)
	assert.NoError(err)
	assert.Equal(originalInput.Name, getAccountOutput.Name.String())
}
