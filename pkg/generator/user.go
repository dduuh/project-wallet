package generator

import (
	"encoding/json"

	"github.com/brianvoe/gofakeit"
)

type UserGenerator struct {
	Name   string
	Age    int64
	Gender string
}

func GenerateUser() ([]byte, error) {
	fakeUser := &UserGenerator{
		Name:   gofakeit.BeerName(),
		Age:    gofakeit.Int64(),
		Gender: gofakeit.Gender(),
	}

	fakeUserBytes, err := json.Marshal(fakeUser)
	if err != nil {
		return nil, err
	}

	return fakeUserBytes, nil
}
