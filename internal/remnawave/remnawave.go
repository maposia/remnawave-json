package remnawave

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"remnawave-json/internal/config"
)

type User struct {
	ShortUuid    string `json:"shortUuid"`
	DaysLeft     int    `json:"daysLeft"`
	TrafficUsed  string `json:"trafficUsed"`
	TrafficLimit string `json:"trafficLimit"`
	Username     string `json:"username"`
	ExpiresAt    string `json:"expiresAt"`
	IsActive     bool   `json:"isActive"`
	UserStatus   string `json:"userStatus"`
}
type SubscriptionResponse struct {
	IsFound         bool                   `json:"isFound"`
	User            User                   `json:"user"`
	Links           []string               `json:"links"`
	SsConfLinks     map[string]interface{} `json:"ssConfLinks"`
	SubscriptionUrl string                 `json:"subscriptionUrl"`
}
type ResponseWrapper struct {
	Response SubscriptionResponse `json:"response"`
}

func GetSubscription(shortUuid string, header string) (*SubscriptionResponse, error) {
	httpReq, err := http.NewRequest(http.MethodGet, config.GetRemnaweveURL()+"/api/sub/"+shortUuid+"/info", nil)
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

func GetUserInfo(shortUuid string, header string) (headers map[string][]string, body string, err error) {
	httpReq, err := http.NewRequest(http.MethodGet, config.GetRemnaweveURL()+"/api/sub/"+shortUuid, nil)
	if err != nil {
		return nil, "", fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("User-Agent", header)
	httpReq.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := config.GetHttpClient().Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("error while getting subscription", slog.String("url", config.GetRemnaweveURL()+resp.Request.URL.String()))
		return nil, "", fmt.Errorf("getting subscription status: %s", resp.Status)
	}

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("reading response body: %w", err)
	}

	body = base64.StdEncoding.EncodeToString(bodyByte)

	filteredHeaders := []string{
		"Profile-Title",
		"Profile-Update-Interval",
		"Subscription-Userinfo",
		"Profile-Web-Page-Url",
		"Content-Disposition",
	}

	resultHeaders := make(map[string][]string)

	for _, header := range filteredHeaders {
		if values, found := resp.Header[header]; found {
			resultHeaders[header] = values
		}
	}

	return resultHeaders, body, nil
}
