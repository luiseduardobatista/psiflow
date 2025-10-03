package domain

import (
	"errors"
	"strings"
)

var validStates = map[string]bool{
	"AC": true, "AL": true, "AP": true, "AM": true, "BA": true, "CE": true, "DF": true,
	"ES": true, "GO": true, "MA": true, "MT": true, "MS": true, "MG": true, "PA": true,
	"PB": true, "PR": true, "PE": true, "PI": true, "RJ": true, "RN": true, "RS": true,
	"RO": true, "RR": true, "SC": true, "SP": true, "SE": true, "TO": true,
}

type AddressState string

func NewAddressState(state string) (*AddressState, error) {
	s := strings.TrimSpace(strings.ToUpper(state))
	if !validStates[s] {
		return nil, errors.New("invalid state abbreviation")
	}
	addressState := AddressState(s)
	return &addressState, nil
}

func (s AddressState) String() string {
	return string(s)
}
