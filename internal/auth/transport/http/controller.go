package http

import (
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

const refreshTokenCookie = "refresh_token"

type AuthController struct {
	Service *service.AuthService
}

func (a *AuthController) RegisterRoutes(s *fuego.Server) {
	m := shared.ParseTokenMiddlewareBuilder{
		Service: a.Service.AccessTokenSvc,
	}.Middleware

	g := fuego.Group(s, "/auth")
	fuego.Post(g, "/login", a.Login)
	fuego.Get(g, "/refresh", a.Refresh)
	fuego.Get(g, "/me", a.Me, option.Middleware(m))
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (a *AuthController) Login(c fuego.ContextWithBody[LoginRequest]) (string, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	tokenPair, err := a.Service.Login(c.Context(), body.Email, body.Password)
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

func (a *AuthController) Refresh(c fuego.ContextNoBody) (string, error) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		return "", fuego.UnauthorizedError{
			Err:    err,
			Detail: "refresh token not found",
		}
	}

	return a.Service.Refresh(c.Context(), refreshToken.Value)
}

func (a *AuthController) Me(c fuego.ContextNoBody) (*service.Account, error) {
	id := c.Context().Value(shared.AccountIDKey).(uuid.UUID)

	return a.Service.Me(c.Context(), id)
}
