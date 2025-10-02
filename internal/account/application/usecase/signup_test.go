package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/luiseduardobatista/psiflow/internal/account/application/usecase"
)

func TestShouldCreateValidAccount(t *testing.T) {
	assert := assert.New(t)
	inputSignup := usecase.SignupInput{
		Name:     "John Doe",
		Email:    fmt.Sprintf("user%d@example.com", time.Now().UnixNano()),
		Password: "123",
		Phone:    "+5511999999999",
	}
	signup := usecase.Signup{}
	signupOutput, err := signup.Execute(inputSignup)
	if err != nil {
		assert.Nil(signupOutput, "signupOutput should not return error")
	}
	assert.NotEmptyf(signupOutput.AccountID, "AccountID should not be empty")
}
