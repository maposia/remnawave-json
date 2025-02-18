package main

import (
	"log/slog"
	"os"
	"os/signal"
	"remnawawe-json/internal/app"
	"remnawawe-json/internal/config"
	"remnawawe-json/internal/service"
	"remnawawe-json/remnawawe"
	"syscall"
)

func main() {
	config.InitConfig()

	conf := config.GetConfig()
	remnawawePanel := remnawawe.NewPanel(conf.RemnaweweURL, conf.RemnawaweToken)

	go app.Start(service.NewService(remnawawePanel))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down...")

	app.Stop()

	slog.Info("Gracefully stopped.")
}
