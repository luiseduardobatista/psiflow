package domain

import (
	"errors"
	"regexp"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z]+ [a-zA-Z\s]*$`)

type Name string

func NewName(name string) (Name, error) {
	if !nameRegex.MatchString(name) {
		return "", errors.New("invalid name: must contain at least a first and last name")
	}
	return Name(name), nil
}

func (n Name) String() string {
	return string(n)
}
