package zipcode

import (
	"errors"
	"unicode"
)

type Zipcode struct {
	Zipcode string
}

func NewZipcode(zipcode string) (Zipcode, error) {
	res := Zipcode{
		zipcode,
	}
	err := res.IsValid()
	if err != nil {
		return Zipcode{}, err
	}
	return res, nil
}

func (h *Zipcode) IsValid() error {
	if len(h.Zipcode) != 8 {
		return errors.New("CEP deve conter 8 dígitos")
	}
	for _, r := range h.Zipcode {
		if !unicode.IsDigit(r) {
			return errors.New("CEP deve conter apenas números")
		}
	}
	return nil
}
