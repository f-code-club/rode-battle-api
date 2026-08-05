package model

import "github.com/google/uuid"

type Role string

const (
	Participant Role = "participant"
	Jury        Role = "jury"
	Admin       Role = "admin"
)

type CreateAccountData struct {
	Email    string
	Name     string
	Password string
	Role     Role
}

type Account struct {
	ID       uuid.UUID
	Email    string
	Name     string
	Password string
	Role     Role
}
