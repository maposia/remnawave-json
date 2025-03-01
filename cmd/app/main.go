package main

import (
	"log/slog"
	"os"
	"os/signal"
	"remnawave-json/internal/app"
	"remnawave-json/internal/config"
	"syscall"
)

func main() {
	config.InitConfig()

	go app.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down...")

	app.Stop()

	slog.Info("Gracefully stopped.")
}
