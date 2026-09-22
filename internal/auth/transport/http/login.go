package http

import (
	"context"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" format:"email" required:"true"`
	Password string `json:"password" minLength:"8" required:"true"`
}

type LoginInput struct {
	Body LoginRequest
}

type LoginOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      string
}

func (s *Server) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	tokenPair, err := s.service.Login(ctx, input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, err
	}

	resp := &LoginOutput{
		SetCookie: http.Cookie{
			Name:  refreshTokenCookie,
			Value: tokenPair.RefreshToken,
		},
		Body: tokenPair.AccessToken,
	}
	return resp, nil
}
