package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"remnawave-json/internal/config"
	"remnawave-json/internal/transport/rest"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var server *http.Server

func Start() {
	r := mux.NewRouter()

	r.Use(httpsAndProxyMiddleware)

	r.HandleFunc("/{shortUuid}", userAgentRouter()).Methods(http.MethodGet)
	r.HandleFunc("/{shortUuid}/v2ray-json", v2rayJson()).Methods(http.MethodGet)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("/app/templates/assets"))))

	server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.GetAppHost(), config.GetAppPort()),
		Handler: r,
	}

	slog.Info("Starting server on http://" + server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Error while starting server")
		panic(err)
	}
}

func v2rayJson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest.V2rayJson(w, r)
		return
	}
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error during server shutdown", "error", err)
		if err = server.Close(); err != nil {
			slog.Error("Error during server shutdown", "error", err)
			panic(err)
		}
	}
}

func httpsAndProxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.GetAppHost() == "localhost" {
			next.ServeHTTP(w, r)
			return
		}
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		xForwardedProto := r.Header.Get("X-Forwarded-Proto")

		if xForwardedFor == "" || xForwardedProto != "https" {
			slog.Error("Reverse proxy and HTTPS are required.")
			http.Error(w, "Reverse proxy and HTTPS are required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func userAgentRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if isBrowser(userAgent) {
			rest.WebPage(w, r)
			return
		}
		if strings.Contains(userAgent, "Streisand") {
			rest.V2rayJson(w, r)
			return
		}

		if strings.Contains(userAgent, "Happ") && config.IsHappJsonEnabled() {
			rest.HappJson(w, r)
			return
		}

		rest.Direct(w, r)
	}
}

var browserKeywords = [...]string{"Mozilla", "Chrome", "Safari", "Firefox", "Opera", "Edge", "TelegramBot"}

func isBrowser(userAgent string) bool {
	for _, keyword := range browserKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}
	return false
}
