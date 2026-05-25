package http

import "github.com/go-fuego/fuego"

func (s *Server) TestHandler(c fuego.ContextNoBody) (string, error) {
	return "test successful!", nil
}
