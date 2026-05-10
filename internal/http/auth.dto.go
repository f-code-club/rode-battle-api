package http

import (
	"context"

	"github.com/go-fuego/fuego"
)

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	Name        string `json:"name" validate:"required"`
	School      string `json:"school" validate:"required"`
	StudentId   string `json:"student_id" validate:"required,alphanum"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric,min=10,max=11"`
}

var _ fuego.InTransformer = (*RegisterRequest)(nil)

func (r *RegisterRequest) InTransform(ctx context.Context) error {
	if !r.isValidPassword() {
		return ErrInvalidPassword
	}

	if !r.isValidName() {
		return ErrInvalidName
	}

	if !r.isValidSchool() {
		return ErrInvalidSchool
	}

	return nil
}
