package shared

import "net/smtp"

type EmailService struct {
	Username   string
	Password   string
	HostName   string
	PortNumber string
}

func NewEmailService(username, password, hostName string, portNumber string) EmailService {
	return EmailService{username, password, hostName, portNumber}
}

func (s *EmailService) SendSingleMail(
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
