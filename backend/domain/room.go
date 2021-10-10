package domain

import "math/rand"

type Room struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}

func NewRoom(name string) *Room {
	return &Room{
		Id:   int32(rand.Int31()),
		Name: name,
	}
}
