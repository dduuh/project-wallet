package generator

import (
	"encoding/json"
	"fmt"

	"github.com/brianvoe/gofakeit"
)

type UserGenerator struct {
	Name   string `json:"name"`
	Age    int64  `json:"age"`
	Gender string `json:"gender"`
}

func GenerateUser() ([]byte, error) {
	fakeUser := &UserGenerator{
		Name:   gofakeit.BeerName(),
		Age:    gofakeit.Int64(),
		Gender: gofakeit.Gender(),
	}

	fakeUserBytes, err := json.Marshal(fakeUser)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UserGenerator: %w", err)
	}

	return fakeUserBytes, nil
}
