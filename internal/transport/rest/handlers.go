package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"remnawawe-json/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) V2rayJson(w http.ResponseWriter, r *http.Request) {
	shortUuid := mux.Vars(r)["shortUuid"]

	jsonData, headers, err := h.service.GenerateJson(shortUuid)
	if err != nil {
		slog.Error("Get Json Error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	for key, values := range headers {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		return
	}
}
