package http

import (
	"github.com/f-code-club/rode-battle-api/internal/mail/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/go-fuego/fuego"
)

type Server struct {
	Service        service.Service
	accessTokenSvc *shared.TokenService
}

func NewServer(
	cfg *shared.Config,
	accessTokenSvc *shared.TokenService,
) Server {
	service := service.New(cfg.MailerUsername, cfg.MailerPassword, cfg.MailerHostName, cfg.MailerPortNumber)

	return Server{service, accessTokenSvc}
}

func (s *Server) SendSingleMailRoute(f *fuego.Server) {
	m := fuego.Group(f, "/mail")

	fuego.Post(m, "/send-single", s.SendSingleMail)
}
