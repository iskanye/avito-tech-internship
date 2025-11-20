package server

import "github.com/iskanye/avito-tech-internship/pkg/api"

type Server struct{}

// Проверка на реализацию всех методов
var _ api.StrictServerInterface = (*Server)(nil)

func NewServer() Server {
	return Server{}
}
