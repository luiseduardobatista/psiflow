package domain

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrInvalidZipCodeFormat        = errors.New("zip code must have 8 digits")
	ErrInvalidZipCodeAllDigitsSame = errors.New("zip code with all same digits is not allowed")
)
var nonDigitRegex = regexp.MustCompile("[^0-9]+")

type ZipCode string

func NewZipCode(zipCode string) (*ZipCode, error) {
	cleaned := nonDigitRegex.ReplaceAllString(zipCode, "")
	if len(cleaned) != 8 {
		return nil, ErrInvalidZipCodeFormat
	}
	firstDigit := cleaned[0]
	isAllSame := true
	for i := 1; i < len(cleaned); i++ {
		if cleaned[i] != firstDigit {
			isAllSame = false
			break
		}
	}
	if isAllSame {
		return nil, ErrInvalidZipCodeAllDigitsSame
	}
	zc := ZipCode(cleaned)
	return &zc, nil
}

func (z ZipCode) String() string {
	return string(z)
}

func (z ZipCode) Formatted() string {
	str := z.String()
	if len(str) != 8 {
		return str
	}
	return fmt.Sprintf("%s-%s", str[0:5], str[5:8])
}
