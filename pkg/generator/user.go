package generator

import "github.com/brianvoe/gofakeit"

type UserGenerator struct {
	Name   string
	Age    int32
	Gender string
}

func NewUserGenerator() *UserGenerator {
	return &UserGenerator{
		Name:   gofakeit.BeerName(),
		Age:    gofakeit.Int32(),
		Gender: gofakeit.Gender(),
	}
}
