package http

import (
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (s *Server) RegisterRoutes(f *fuego.Server) {
	fuego.Get(f, "/", s.PingHandler)

	auth := fuego.Group(f, "/auth")
	fuego.Post(auth, "/register", s.RegisterHandler)
	fuego.Post(auth, "/login", s.LoginHandler)

	test := fuego.Group(f, "/test",
		option.Middleware(s.authMiddleware),
	)
	fuego.Get(test, "/middleware", s.TestHandler)
}

func (s *Server) OpenAPIHandler(specURL string) http.Handler {
	return httpSwagger.Handler(
		httpSwagger.Layout(httpSwagger.StandaloneLayout),
		httpSwagger.PersistAuthorization(true),
		httpSwagger.URL(specURL),
	)
}

func (s *Server) PingHandler(c fuego.ContextNoBody) (string, error) {
	return "pong", nil
}
