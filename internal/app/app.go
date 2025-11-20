package app

import (
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/avito-tech-internship/internal/api"
	"github.com/iskanye/avito-tech-internship/internal/config"
	"github.com/iskanye/avito-tech-internship/internal/repositories"
	"github.com/iskanye/avito-tech-internship/internal/server"
)

type App struct {
	e   *gin.Engine
	s   *repositories.Storage
	cfg *config.Config
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

func (a *App) MustRun() {
	if err := a.e.Run(address(a.cfg.Host, a.cfg.Port)); err != nil {
		panic(err)
	}
}

func (a *App) GracefulStop() {
	a.s.Stop()
}

func address(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
