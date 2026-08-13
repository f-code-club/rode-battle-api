package http

import (
	"github.com/f-code-club/rode-battle-api/internal/account/service"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

type Role string

const (
	Participant Role = "participant"
	Jury        Role = "jury"
	Admin       Role = "admin"
)

type GenerateRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
	Role  Role   `json:"role" validate:"required"`
}

func (s *Server) Generate(c fuego.ContextWithBody[GenerateRequest]) (uuid.UUID, error) {
	body, err := c.Body()
	if err != nil {
		return uuid.Nil, err
	}

	return s.service.Generate(c.Context(), body.Email, body.Name, service.Role(body.Role))
}
