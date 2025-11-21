package app

import (
	"log/slog"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/avito-tech-internship/internal/config"
	"github.com/iskanye/avito-tech-internship/internal/repositories"
	"github.com/iskanye/avito-tech-internship/internal/server"
	"github.com/iskanye/avito-tech-internship/internal/service/prassignment"
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
	storage, err := repositories.New(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.MaxConns,
	)
	if err != nil {
		panic(err)
	}

	prAssignment := prassignment.New(
		log,
		storage, storage, storage,
		storage, storage,
	)
	server.Register(engine, prAssignment)

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
