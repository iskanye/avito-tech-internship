package app

import (
	"log/slog"
	"net"
	"strconv"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
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
		trmpgx.DefaultCtxGetter,
	)
	if err != nil {
		panic(err)
	}

	txManager := manager.Must(trmpgx.NewDefaultFactory(storage.GetPool()))

	// Это страшно
	prAssignment := prassignment.New(
		log,
		txManager,
		storage, storage, storage,
		storage, storage,
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
