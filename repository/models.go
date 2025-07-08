package repository

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Todo struct {
	ID   uuid.UUID `json:"id"`
	Text string    `json:"text"`
	Done bool      `json:"done"`
	User *User     `json:"user"`
}
