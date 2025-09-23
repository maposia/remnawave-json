package rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"remnawave-json/internal/config"
	"remnawave-json/internal/remnawave"

	"github.com/gorilla/mux"
)

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

	fmt.Println(sub)

	wrapped := struct {
		Response interface{} `json:"response"`
	}{
		Response: sub,
	}

	jsonData, err := json.Marshal(wrapped)
	if err != nil {
		http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
		return
	}

	panelDataB64 := base64.StdEncoding.EncodeToString(jsonData)

	data := struct {
		MetaTitle, MetaDescription, PanelData string
	}{
		PanelData: panelDataB64, MetaTitle: config.GetMetaTitle(), MetaDescription: config.GetMetaDescription(),
	}

	err = config.GetWebPageTemplate().Execute(w, data)
	if err != nil {
		slog.Error("Execute Json Error", err)
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
	if err != nil {
		http.Error(w, "failed to read response body", http.StatusInternalServerError)
		return
	}

	data, err := DecodeJSON(body)
	if err != nil {
		log.Printf("JSON parse error: %v", err)
	} else {
		//if config.GetRuHostName() != "" {
		//	rawSub, err := remnawave.GetRawSubscription(shortUuid, r.Header.Get("User-Agent"))
		//	if err != nil {
		//		log.Printf("JSON parse error: %v", err)
		//	}
		//
		//	ruHost := findRawHostByRemark(rawSub, config.GetRuHostName())
		//	if ruHost != nil {
		//		UpdateRuOutbound(data, ruHost)
		//	}
		//}
	}

	if _, exists := config.GetExceptRuRulesUsers()[shortUuid]; exists {
		data = CleanRURules(data)
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

func CleanRURules(data interface{}) interface{} {
	arr, ok := data.([]interface{})
	if !ok {
		return data
	}

	for i, v := range arr {
		obj, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		routing, ok := obj["routing"].(map[string]interface{})
		if !ok {
			continue
		}

		if rules, ok := routing["rules"].([]interface{}); ok {
			filtered := make([]interface{}, 0, len(rules))
			for _, r := range rules {
				if rule, ok := r.(map[string]interface{}); ok {
					if tag, ok := rule["outboundTag"].(string); ok && tag == "RU" {
						continue
					}
				}
				filtered = append(filtered, r)
			}
			routing["rules"] = filtered
		}

		arr[i] = obj
	}

	return arr
}

func HappJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("routing", config.GetHappRouting())
	r.Header.Set("User-Agent", r.Header.Get("User-Agent"))
	V2rayJson(w, r)
}

func BalancerConfig(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("User-Agent", r.Header.Get("User-Agent"))
	BalancerJson(w, r)
}

func BalancerJson(w http.ResponseWriter, r *http.Request) {
	shortUuid := mux.Vars(r)["shortUuid"]

	rawData, err := remnawave.GetRawSubscription(shortUuid, r.Header.Get("User-Agent"))
	if err != nil {
		log.Printf("Failed to get raw subscription: %v", err)
		return
	}

	xrayConfig, err := remnawave.ConvertToXrayConfig(rawData)
	if err != nil {
		return
	}

	data, err := DecodeJSON(xrayConfig)
	if err != nil {
		log.Printf("JSON parse error: %v", err)
	}

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
	//if _, exists := config.GetExceptRuRulesUsers()[shortUuid]; exists {
	//	data = CleanRURules(data)
	//}
	//for key, values := range resp.Header {
	//	for _, value := range values {
	//		w.Header().Add(key, value)
	//	}
	//}

	w.Header().Set("routing", "happ://routing/add/eyJOYW1lIjoiIiwiR2xvYmFsUHJveHkiOiJ0cnVlIiwiUmVtb3RlRE5TVHlwZSI6IkRvSCIsIlJlbW90ZUROU0RvbWFpbiI6IiIsIlJlbW90ZUROU0lQIjoiIiwiRG9tZXN0aWNETlNUeXBlIjoiRG9VIiwiRG9tZXN0aWNETlNEb21haW4iOiIiLCJEb21lc3RpY0ROU0lQIjoiIiwiR2VvaXB1cmwiOiJodHRwczovL2dpdGh1Yi5jb20vZnJheVpWL3NpbXBsZS1ydS1nZW9pcC9yZWxlYXNlcy9sYXRlc3QvZG93bmxvYWQvZ2VvaXAuZGF0IiwiR2Vvc2l0ZXVybCI6Imh0dHBzOi8vZ2l0aHViLmNvbS9mcmF5WlYvc2ltcGxlLXJ1LWdlb3NpdGUvcmVsZWFzZXMvbGF0ZXN0L2Rvd25sb2FkL2dlb3NpdGUuZGF0IiwiTGFzdFVwZGF0ZWQiOiIiLCJEbnNIb3N0cyI6e30sIkRpcmVjdFNpdGVzIjpbXSwiRGlyZWN0SXAiOltdLCJQcm94eVNpdGVzIjpbXSwiUHJveHlJcCI6W10sIkJsb2NrU2l0ZXMiOltdLCJCbG9ja0lwIjpbXSwiRG9tYWluU3RyYXRlZ3kiOiJJUElmTm9uTWF0Y2giLCJGYWtlRE5TIjoiZmFsc2UiLCJVc2VDaHVua0ZpbGVzIjoidHJ1ZSJ9")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
	}
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
