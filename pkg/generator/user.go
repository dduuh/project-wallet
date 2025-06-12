package generator

import "github.com/brianvoe/gofakeit"

type UserGenerator struct {
	Name   string
	Age    uint8
	Gender string
}

func NewUserGenerator() *UserGenerator {
	return &UserGenerator{
		Name: gofakeit.BeerName(),
		Age:  uint8(gofakeit.Int8()),
		Gender: gofakeit.Gender(),
	}
}