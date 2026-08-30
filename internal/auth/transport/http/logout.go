package http

import (
	"net/http"

	"github.com/go-fuego/fuego"
)

func (s *Server) Logout(c fuego.ContextNoBody) (string, error) {
	c.SetCookie(http.Cookie{
		Name:   refreshTokenCookie,
		Value:  "",
		MaxAge: -1,
	})

	return "Logged out successfully", nil
}
