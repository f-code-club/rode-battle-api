package service

type Service struct {
	Username   string
	Password   string
	HostName   string
	PortNumber string
}

func New(username, password, hostName string, portNumber string) Service {
	return Service{username, password, hostName, portNumber}
}
