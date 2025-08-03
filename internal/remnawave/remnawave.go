package remnawave

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"remnawave-json/internal/config"
)

type SubscriptionResponse struct {
	IsFound         bool                   `json:"isFound"`
	User            User                   `json:"user"`
	Links           []string               `json:"links"`
	SsConfLinks     map[string]interface{} `json:"ssConfLinks"`
	SubscriptionUrl string                 `json:"subscriptionUrl"`
}

type Response struct {
	User            User      `json:"user"`
	SubscriptionURL string    `json:"subscriptionUrl"`
	IsHwidLimited   bool      `json:"isHwidLimited"`
	RawHosts        []RawHost `json:"rawHosts"`
	Headers         Headers   `json:"headers"`
}

type User struct {
	ShortUUID            string `json:"shortUuid"`
	DaysLeft             int    `json:"daysLeft"`
	TrafficUsed          string `json:"trafficUsed"`
	TrafficLimit         string `json:"trafficLimit"`
	Username             string `json:"username"`
	ExpiresAt            string `json:"expiresAt"`
	IsActive             bool   `json:"isActive"`
	UserStatus           string `json:"userStatus"`
	TrafficLimitStrategy string `json:"trafficLimitStrategy"`
}

type RawHost struct {
	Remark        string      `json:"remark"`
	Address       string      `json:"address"`
	Port          int         `json:"port"`
	Protocol      string      `json:"protocol"`
	Path          string      `json:"path"`
	Host          string      `json:"host"`
	TLS           string      `json:"tls"`
	SNI           string      `json:"sni"`
	ALPN          string      `json:"alpn"`
	PublicKey     string      `json:"publicKey"`
	Fingerprint   string      `json:"fingerprint"`
	ShortID       string      `json:"shortId"`
	SpiderX       string      `json:"spiderX"`
	Network       string      `json:"network"`
	Password      Passwords   `json:"password"`
	MuxParams     MuxParams   `json:"muxParams"`
	SockoptParams interface{} `json:"sockoptParams"` // null в JSON
	Flow          string      `json:"flow"`
}

type Passwords struct {
	TrojanPassword string `json:"trojanPassword"`
	VlessPassword  string `json:"vlessPassword"`
	SSPassword     string `json:"ssPassword"`
}

type MuxParams struct {
	Enabled         bool   `json:"enabled"`
	Concurrency     int    `json:"concurrency"`
	XudpConcurrency int    `json:"xudpConcurrency"`
	XudpProxyUDP443 string `json:"xudpProxyUDP443"`
}

type Headers struct {
	ContentDisposition     string `json:"content-disposition"`
	SupportURL             string `json:"support-url"`
	ProfileTitle           string `json:"profile-title"`
	ProfileUpdateInterval  string `json:"profile-update-interval"`
	SubscriptionUserInfo   string `json:"subscription-userinfo"`
	Announce               string `json:"announce"`
	ProfileWebPageURL      string `json:"profile-web-page-url"`
	SubscriptionRefillDate string `json:"subscription-refill-date"`
	HideSettings           string `json:"hide-settings"`
}

type ResponseWrapper struct {
	Response SubscriptionResponse `json:"response"`
}

func GetSubscription(shortUuid string, header string) (*SubscriptionResponse, error) {
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/sub/%s/info", config.GetRemnaweveURL(), shortUuid), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", header)

	resp, err := config.GetHttpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("error while getting subscription", slog.String("url", config.GetRemnaweveURL()+resp.Request.URL.String()))
		return nil, fmt.Errorf("getting subscription status: %s", resp.Status)
	}

	response := &ResponseWrapper{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return &response.Response, nil
}

func GetRawSubscription(shortUuid string, header string) (*Response, error) {
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/sub/%s/raw", config.GetRemnaweveURL(), shortUuid), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", header)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetRemnawaveToken()))

	resp, err := config.GetHttpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("error while getting subscription", slog.String("url", config.GetRemnaweveURL()+resp.Request.URL.String()))
		return nil, fmt.Errorf("getting subscription status: %s", resp.Status)
	}

	response := &Response{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return response, nil
}
