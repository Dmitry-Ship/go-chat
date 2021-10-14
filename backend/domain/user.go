package domain

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

type User struct {
	Id     uuid.UUID `json:"id"`
	Avatar string    `json:"avatar"`
	Name   string    `json:"name"`
}

func getRandomWord() string {
	randomLetters := ""
	for i := 0; i < 6; i++ {
		randomLetters += string(rune(rand.Intn(26) + 97))

		if i == 0 {
			randomLetters = strings.ToUpper(randomLetters)
		}

	}

	return randomLetters
}

func NewUser() *User {
	firstName := getRandomWord()
	lastName := getRandomWord()

	return &User{
		Id:     uuid.New(),
		Avatar: string(firstName[0]) + string(lastName[0]),
		Name:   firstName + " " + lastName,
	}
}
