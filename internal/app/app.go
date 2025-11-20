package app

import (
	"log/slog"
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
	log *slog.Logger
	cfg *config.Config
}

func New(
	engine *gin.Engine,
	log *slog.Logger,
	cfg *config.Config,
) App {
	server := server.NewServer()

	api.RegisterHandlers(engine, api.NewStrictHandler(
		server,
		[]api.StrictMiddlewareFunc{},
	))

	return App{
		e:   engine,
		log: log,
		cfg: cfg,
	}
}

func (a App) MustRun() {
	if err := a.e.Run(address(a.cfg.Host, a.cfg.Port)); err != nil {
		panic(err)
	}
}

func (a App) GracefulStop() {
	a.s.Stop()
	a.log.Info("Gracefully stopped")
}

func address(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
