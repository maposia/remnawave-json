package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log/slog"
	"net/http"
	"remnawave-json/internal/config"
	"remnawave-json/internal/remnawave"
	"remnawave-json/internal/service"
	"time"
)

func V2rayJson(w http.ResponseWriter, r *http.Request) {
	shortUuid := mux.Vars(r)["shortUuid"]
	header := r.Header.Get("User-Agent")

	jsonData, headers, err := service.GenerateJson(shortUuid, header)
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

func V2ray(w http.ResponseWriter, r *http.Request) {
	shortUuid := mux.Vars(r)["shortUuid"]

	proxyURL := config.GetRemnaweveURL() + "/api/sub/" + shortUuid
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

	httpReq.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := config.GetHttpClient().Do(httpReq)
	if err != nil {
		http.Error(w, "failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

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

func Direct(w http.ResponseWriter, r *http.Request) {
	shortUuid := mux.Vars(r)["shortUuid"]

	proxyURL := config.GetRemnaweveURL() + "/api/sub/" + shortUuid
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

	resp, err := config.GetHttpClient().Do(httpReq)
	if err != nil {
		http.Error(w, "failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

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

type WebPageUser struct {
	Status          string
	ExpireFormatted string
	UsedTraffic     string
	DataLimit       string
	SubscriptionURL string
	ResetStrategy   string
	Username        string
}

func WebPage(w http.ResponseWriter, r *http.Request) {
	shortUuid := mux.Vars(r)["shortUuid"]
	header := r.Header.Get("User-Agent")
	sub, err := remnawave.GetSubscription(shortUuid, header)
	if err != nil {
		slog.Error("Get Json Error", err)
		http.Error(w, "Ошибка получения подписки", http.StatusInternalServerError)
		return
	}

	var expireFormatted string
	if sub.User.ExpiresAt != "" {
		expireTime, err := time.Parse(time.RFC3339, sub.User.ExpiresAt)
		if err != nil {
			slog.Error("Invalid date format", err)
			expireFormatted = ""
		} else {
			expireFormatted = expireTime.Format(time.RFC3339)
		}
	}

	user := WebPageUser{
		Status:          sub.User.UserStatus,
		ExpireFormatted: expireFormatted,
		UsedTraffic:     sub.User.TrafficUsed,
		DataLimit:       sub.User.TrafficLimit,
		SubscriptionURL: sub.SubscriptionUrl,
		Username:        sub.User.Username,
	}

	err = config.GetWebPageTemplate().Execute(w, user)
	if err != nil {
		http.Error(w, "Ошибка заполнения шаблона", http.StatusInternalServerError)
	}
}

func Streisand(w http.ResponseWriter, r *http.Request) {
	shortUuid := mux.Vars(r)["shortUuid"]

	proxyURL := config.GetRemnaweveURL() + "/api/sub/" + shortUuid + "/v2ray-json"
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

	resp, err := config.GetHttpClient().Do(httpReq)
	if err != nil {
		http.Error(w, "failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

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
