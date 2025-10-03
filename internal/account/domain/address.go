package domain

type Address struct {
	Street       string
	Neighborhood string
	Number       AddressNumber
	City         string
	State        AddressState
	ZipCode      ZipCode
}

func NewAddress(street string, neighborhood string, number int, city string, state string, zipCode string) (*Address, error) {
	addressNumber, err := NewAddressNumber(number)
	if err != nil {
		return nil, err
	}
	addressState, err := NewAddressState(state)
	if err != nil {
		return nil, err
	}
	addressZipCode, err := NewZipCode(zipCode)
	if err != nil {
		return nil, err
	}
	return &Address{
		Street:       street,
		Neighborhood: neighborhood,
		Number:       *addressNumber,
		City:         city,
		State:        *addressState,
		ZipCode:      *addressZipCode,
	}, nil
}
