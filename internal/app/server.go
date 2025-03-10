package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"regexp"
	"remnawave-json/internal/config"
	"remnawave-json/internal/transport/rest"
	"strings"
	"time"
)

var server *http.Server

func Start() {
	r := mux.NewRouter()

	r.Use(httpsAndProxyMiddleware)

	r.HandleFunc("/{shortUuid}/v2ray-json", rest.V2rayJson).Methods("GET")
	r.HandleFunc("/{shortUuid}", userAgentRouter()).Methods("GET")

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

var v2rayNRegex = regexp.MustCompile(`^v2rayN/(\d+\.\d+)`)
var v2rayNGRegex = regexp.MustCompile(`^v2rayNG/(\d+\.\d+\.\d+)`)
var streisandRegex = regexp.MustCompile(`^[Ss]treisand`)
var happRegex = regexp.MustCompile(`^Happ/`)
var ktorClientRegex = regexp.MustCompile(`^ktor-client`)
var v2boxRegex = regexp.MustCompile(`^V2Box`)

func userAgentRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		switch {
		case v2rayNRegex.MatchString(userAgent):
			version := v2rayNRegex.FindStringSubmatch(userAgent)[1]
			if compareVersions(version, "6.40") >= 0 {
				rest.V2rayJson(w, r)
			} else {
				rest.Direct(w, r)
			}

		case v2rayNGRegex.MatchString(userAgent):
			version := v2rayNGRegex.FindStringSubmatch(userAgent)[1]
			if compareVersions(version, "1.8.29") >= 0 {
				rest.V2rayJson(w, r)
			} else {
				rest.Direct(w, r)
			}

		case streisandRegex.MatchString(userAgent):
			rest.V2rayJson(w, r)

		case happRegex.MatchString(userAgent):
			if config.GetHappAnnouncements() != "" {
				w.Header().Set("announce", config.GetHappAnnouncements())
			}
			if config.GetHappRouting() != "" {
				w.Header().Set("routing", config.GetHappRouting())
			}
			if config.IsHappJsonEnabled() {
				rest.V2rayJson(w, r)
			} else {
				rest.Direct(w, r)
			}

		case ktorClientRegex.MatchString(userAgent):
			rest.V2rayJson(w, r)

		case v2boxRegex.MatchString(userAgent):
			rest.V2rayJson(w, r)

		default:
			if isBrowser(userAgent) {
				rest.WebPage(w, r)
				return
			}
			rest.Direct(w, r)
		}
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

func compareVersions(version1, version2 string) int {
	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")

	for i := 0; i < len(v1) && i < len(v2); i++ {
		if v1[i] > v2[i] {
			return 1
		} else if v1[i] < v2[i] {
			return -1
		}
	}
	return 0
}
