package http

import (
	"net/http"

	"github.com/go-fuego/fuego"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (s *Server) Login(c fuego.ContextWithBody[LoginRequest]) (string, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	tokenPair, err := s.service.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return "", err
	}

	c.SetCookie(http.Cookie{
		Name:     refreshTokenCookie,
		Value:    tokenPair.RefreshToken,
		SameSite: http.SameSiteNoneMode,
	})
	return tokenPair.AccessToken, nil
}
