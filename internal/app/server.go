package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"remnawawe-json/internal/config"
	"remnawawe-json/internal/service"
	"remnawawe-json/internal/transport/rest"
	"time"
)

var Server *http.Server

func Start(service *service.Service) {
	handler := rest.NewHandler(service)

	r := mux.NewRouter()

	r.HandleFunc("/{shortUuid}/v2ray-json", handler.V2rayJson).Methods("GET")
	r.HandleFunc("/{shortUuid}", handler.Direct).Methods("GET")

	addr := fmt.Sprintf("%s:%s", "localhost", config.GetConfig().AppPort)

	Server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	slog.Info("Starting server on http://" + addr)
	if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Error while starting server")
		panic(err)
	}
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Server.Shutdown(ctx); err != nil {
		slog.Error("Error during server shutdown", "error", err)
		if err = Server.Close(); err != nil {
			slog.Error("Error during server shutdown", "error", err)
			panic(err)
		}
	}

}
