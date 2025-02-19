package rest

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log/slog"
	"net/http"
	"remnawawe-json/internal/service"
	"strings"
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

func (h *Handler) Direct(w http.ResponseWriter, r *http.Request) {
	shortUuid := mux.Vars(r)["shortUuid"]

	proxyURL := h.service.Panel.BaseURL + "/api/sub/" + shortUuid
	httpReq, err := http.NewRequest(r.Method, proxyURL, r.Body)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	resp, err := h.service.Panel.Client.Do(httpReq)
	if err != nil {
		http.Error(w, "failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Копируем заголовки ответа
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "failed to copy response body", http.StatusInternalServerError)
	}
}

func decodeBase64Title(encoded string) string {
	trimmed := strings.TrimPrefix(encoded, "base64:")
	decodedBytes, _ := base64.StdEncoding.DecodeString(trimmed)
	return string(decodedBytes)
}
