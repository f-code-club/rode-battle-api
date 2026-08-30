package http

import (
	"net/http"

	"github.com/go-fuego/fuego"
)

func (s *Server) Logout(c fuego.ContextNoBody) (string, error) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		return "", fuego.UnauthorizedError{
			Err:    err,
			Detail: "refresh token not found",
		}
	}

	err = s.service.Logout(c.Context(), refreshToken.Value)
	if err != nil {
		return "", fuego.UnauthorizedError{
			Err:    err,
			Detail: "invalid refresh token",
		}
	}

	c.SetCookie(http.Cookie{
		Name:   refreshTokenCookie,
		Value:  "",
		MaxAge: -1,
	})

	return "Logged out successfully", nil
}
