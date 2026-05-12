package http

import (
	"context"
	"strings"

	"github.com/go-fuego/fuego"
)

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,strong_password"`
	Name        string `json:"name" validate:"required,human_name"`
	School      string `json:"school" validate:"required,human_name"`
	StudentId   string `json:"student_id" validate:"required,alphanum"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric,min=10,max=11"`
}

var _ fuego.InTransformer = (*RegisterRequest)(nil)

func (r *RegisterRequest) InTransform(ctx context.Context) error {
	r.Email = strings.TrimSpace(
		strings.ToLower(r.Email),
	)

	r.Name = strings.TrimSpace(r.Name)

	r.School = strings.TrimSpace(r.School)

	return nil
}
