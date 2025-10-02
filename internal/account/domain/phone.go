package domain

import (
	"errors"
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type Phone string

func NewPhone(phone string) (Phone, error) {
	const defaultRegion = "BR"
	parsedNumber, err := phonenumbers.Parse(phone, defaultRegion)
	if err != nil {
		return "", errors.New("invalid phone number: input contains non-numeric characters or is incorrectly formatted")
	}
	if !phonenumbers.IsValidNumber(parsedNumber) {
		return "", errors.New("invalid phone number: must be a valid Brazilian number, including the area code")
	}
	formattedPhone := phonenumbers.Format(parsedNumber, phonenumbers.E164)
	return Phone(formattedPhone), nil
}

func (p Phone) String() string {
	return string(p)
}

func (p Phone) Formatted() (string, error) {
	parsedNumber, err := phonenumbers.Parse(string(p), "ZZ")
	if err != nil {
		return "", fmt.Errorf("failed to format phone number: could not re-parse the stored E.164 value '%s': %w", p, err)
	}
	return phonenumbers.Format(parsedNumber, phonenumbers.NATIONAL), nil
}
