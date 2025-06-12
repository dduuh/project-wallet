package domain

type User struct {
	Name   string `json:"name"`
	Age    uint8  `json:"age"`
	Gender string `json:"gender"`
}
