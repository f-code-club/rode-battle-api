package mail

import "net/smtp"

type Mailer struct {
	Username   string
	Password   string
	HostName   string
	PortNumber string
}

type SendSingleMailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

func NewMailService(username, password, hostName string, portNumber string) Mailer {
	return Mailer{username, password, hostName, portNumber}
}

func (s *Mailer) SendSingleMail(
	to string,
	subject string,
	body string,
) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.HostName)

	return smtp.SendMail(
		s.HostName+":"+s.PortNumber,
		auth,
		s.Username,
		[]string{to},
		[]byte(
			"From: "+s.Username+"\r\n"+
				"To: "+to+"\r\n"+
				"Subject: "+subject+"\r\n"+
				"\r\n"+
				body,
		),
	)
}
