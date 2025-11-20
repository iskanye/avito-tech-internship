package app

import (
	"github.com/gin-gonic/gin"
	"github.com/iskanye/avito-tech-internship/internal/api"
	"github.com/iskanye/avito-tech-internship/internal/server"
)

type App struct {
	e *gin.Engine
}

func New(
	engine *gin.Engine,
) *App {
	server := server.NewServer()

	api.RegisterHandlers(engine, api.NewStrictHandler(
		server,
		[]api.StrictMiddlewareFunc{},
	))

	return &App{
		e: engine,
	}
}
