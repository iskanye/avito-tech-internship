package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/avito-tech-internship/internal/app"
	"github.com/iskanye/avito-tech-internship/internal/config"
)

func main() {
	cfg := config.MustLoad()
	cfg.LoadEnv()

	app := app.New(gin.Default(), slog.Default(), cfg)

	go func() {
		app.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GracefulStop()
}
