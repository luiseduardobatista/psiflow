package domain

import "errors"

type AuthenticationProvider interface {
	Signup(email Email, password Password) (string, err error)
}

var ErrAuthEmailAlreadyExists = errors.New("email already exists")
