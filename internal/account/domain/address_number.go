package domain

import "errors"

type AddressNumber int

var ErrInvalidAddressNumber = errors.New("number should be bigger than 0")

func NewAddressNumber(number int) (*AddressNumber, error) {
	if number <= 0 {
		return nil, ErrInvalidAddressNumber
	}
	addressNumber := AddressNumber(number)
	return &addressNumber, nil
}

func (a AddressNumber) Int() int {
	return int(a)
}
