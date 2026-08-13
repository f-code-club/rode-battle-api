package http

import "github.com/go-fuego/fuego"

func (s *Server) Refresh(c fuego.ContextNoBody) (string, error) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		return "", fuego.UnauthorizedError{
			Err:    err,
			Detail: "refresh token not found",
		}
	}

	return s.service.Refresh(c.Context(), refreshToken.Value)
}
