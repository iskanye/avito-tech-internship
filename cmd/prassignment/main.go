package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/avito-tech-internship/internal/app"
	"github.com/iskanye/avito-tech-internship/internal/config"
	"github.com/iskanye/avito-tech-internship/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	cfg.LoadEnv()

	e := gin.New()
	e.Use(gin.Recovery())

	log := logger.SetupPrettySlog()

	app := app.New(e, log, cfg)

	go func() {
		app.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GracefulStop()
}
