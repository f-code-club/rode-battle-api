package http

import (
	"context"
	"net/http"
)

type LogoutInput struct{}

type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      string
}

func (s *Server) Logout(ctx context.Context, input *LogoutInput) (*LogoutOutput, error) {
	return &LogoutOutput{
		SetCookie: http.Cookie{
			Name:   refreshTokenCookie,
			Value:  "",
			MaxAge: -1,
		},
		Body: "Logged out successfully",
	}, nil
}
