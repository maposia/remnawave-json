package main

import (
	"log/slog"
	"os"
	"os/signal"
	"remnawave-json/internal/app"
	"remnawave-json/internal/config"
	"remnawave-json/internal/service"
	"remnawave-json/remnawave"
	"syscall"
)

func main() {
	config.InitConfig()

	conf := config.GetConfig()
	remnawavePanel := remnawave.NewPanel(conf.RemnaweveURL)

	go app.Start(service.NewService(remnawavePanel))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down...")

	app.Stop()

	slog.Info("Gracefully stopped.")
}
