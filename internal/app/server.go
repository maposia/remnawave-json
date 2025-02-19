package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"regexp"
	"remnawawe-json/internal/config"
	"remnawawe-json/internal/service"
	"remnawawe-json/internal/transport/rest"
	"strings"
	"time"
)

var Server *http.Server

func Start(service *service.Service) {
	handler := rest.NewHandler(service)

	r := mux.NewRouter()

	r.HandleFunc("/{shortUuid}/v2ray-json", handler.V2rayJson).Methods("GET")
	r.HandleFunc("/{shortUuid}", userAgentRouter(handler)).Methods("GET")

	Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "localhost", config.GetConfig().AppPort),
		Handler: r,
	}

	slog.Info("Starting server on http://" + Server.Addr)
	if err := Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

var v2rayNRegex = regexp.MustCompile(`^v2rayN/(\d+\.\d+)`)
var v2rayNGRegex = regexp.MustCompile(`^v2rayNG/(\d+\.\d+\.\d+)`)
var streisandRegex = regexp.MustCompile(`^[Ss]treisand`)
var happRegex = regexp.MustCompile(`^Happ/(\d+\.\d+\.\d+)`)
var ktorClientRegex = regexp.MustCompile(`^ktor-client`)
var v2boxRegex = regexp.MustCompile(`^V2Box`)

func userAgentRouter(handler *rest.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")

		switch {
		case v2rayNRegex.MatchString(userAgent):
			version := v2rayNRegex.FindStringSubmatch(userAgent)[1]
			if compareVersions(version, "6.40") >= 0 {
				handler.V2rayJson(w, r)
			} else {
				handler.Direct(w, r)
			}

		case v2rayNGRegex.MatchString(userAgent):
			version := v2rayNGRegex.FindStringSubmatch(userAgent)[1]
			if compareVersions(version, "1.8.29") >= 0 {
				handler.V2rayJson(w, r)
			} else {
				handler.Direct(w, r)
			}

		case streisandRegex.MatchString(userAgent):
			handler.V2rayJson(w, r)

		case happRegex.MatchString(userAgent):
			handler.V2rayJson(w, r)

		case ktorClientRegex.MatchString(userAgent):
			handler.V2rayJson(w, r)

		case v2boxRegex.MatchString(userAgent):
			handler.V2rayJson(w, r)

		default:
			handler.Direct(w, r)
		}
	}
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
