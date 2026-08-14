package http

import (
	"fmt"

	"github.com/go-fuego/fuego"
)

type SendSingleMailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

func (s *Server) SendSingleMail(c fuego.ContextWithBody[SendSingleMailRequest]) (string, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	fmt.Printf("Sending mail from %s with password %s\n", s.Service.Username, s.Service.Password)

	err = s.Service.SendSingleMail(
		body.To,
		body.Subject,
		body.Body,
	)
	if err != nil {
		return "", err
	}

	return "mail sent successfully", nil
}
