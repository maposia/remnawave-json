package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"log/slog"
	"net/http"
	"remnawave-json/internal/config"
	"remnawave-json/internal/remnawave"
	"time"
)

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

func V2rayJson(w http.ResponseWriter, r *http.Request) {
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

	body, err := io.ReadAll(resp.Body)
	data, err := DecodeJSON(body)
	if err != nil {
		log.Printf("JSON parse error: %v", err)
	} else {
		if config.GetRuHostName() != "" {
			rawSub, err := remnawave.GetRawSubscription(shortUuid, r.Header.Get("User-Agent"))
			if err != nil {
				log.Printf("JSON parse error: %v", err)
			}

			ruHost := findRawHostByRemark(rawSub, config.GetRuHostName())
			if ruHost != nil {
				UpdateRuOutbound(data, ruHost)
			}
		}
	}
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
	}
}

func HappJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("routing", config.GetHappRouting())
	r.Header.Set("User-Agent", r.Header.Get("User-Agent"))
	V2rayJson(w, r)
}

func DecodeJSON(body []byte) (interface{}, error) {
	var data interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return data, nil
}

func UpdateRuOutbound(data interface{}, host *remnawave.RawHost) {
	if arr, ok := data.([]interface{}); ok {
		for _, v := range arr {
			UpdateRuOutbound(v, host)
		}
		return
	}

	if m, ok := data.(map[string]interface{}); ok {
		if outbounds, exists := m["outbounds"]; exists {
			if obs, ok := outbounds.([]interface{}); ok {
				for _, ob := range obs {
					if obMap, ok := ob.(map[string]interface{}); ok {
						if tag, hasTag := obMap["tag"]; hasTag && tag == config.GetRuOutboundName() {

							if settings, ok := obMap["settings"].(map[string]interface{}); ok {
								if vnextArr, ok := settings["vnext"].([]interface{}); ok && len(vnextArr) > 0 {
									if vnext, ok := vnextArr[0].(map[string]interface{}); ok {
										if usersArr, ok := vnext["users"].([]interface{}); ok && len(usersArr) > 0 {
											if user, ok := usersArr[0].(map[string]interface{}); ok {
												user["id"] = host.Password.VlessPassword
											}
										}
									}
								}
							}

							if streamSettings, ok := obMap["streamSettings"].(map[string]interface{}); ok {
								if realitySettings, ok := streamSettings["realitySettings"].(map[string]interface{}); ok {
									realitySettings["publicKey"] = host.PublicKey
									realitySettings["shortId"] = host.ShortID
								}
							}
						}
					}
				}
			}
		}
		for _, v := range m {
			UpdateRuOutbound(v, host)
		}
	}
}

func findRawHostByRemark(resp *remnawave.Response, remark string) *remnawave.RawHost {
	for i, host := range resp.RawHosts {
		if host.Remark == remark {
			return &resp.RawHosts[i]
		}
	}
	return nil
}
