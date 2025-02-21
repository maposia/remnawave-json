package remnawave

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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

type Panel struct {
	Client  *http.Client
	BaseURL string
}

func NewPanel(baseURL string) *Panel {
	return &Panel{
		Client: &http.Client{
			Transport: &http.Transport{},
		},
		BaseURL: baseURL,
	}
}

func (p *Panel) GetSubscription(shortUuid string, header string) (*SubscriptionResponse, error) {
	httpReq, err := http.NewRequest(http.MethodGet, p.BaseURL+"/api/sub/"+shortUuid+"/info", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", header)

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("error while getting subscription", slog.String("url", p.BaseURL+resp.Request.URL.String()))
		return nil, fmt.Errorf("getting subscription status: %s", resp.Status)
	}

	response := &ResponseWrapper{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return &response.Response, nil
}

func (p *Panel) GetUserInfo(shortUuid string, header string) (map[string][]string, error) {
	httpReq, err := http.NewRequest(http.MethodGet, p.BaseURL+"/api/sub/"+shortUuid, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("User-Agent", header)

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("error while getting subscription", slog.String("url", p.BaseURL+resp.Request.URL.String()))
		return nil, fmt.Errorf("getting subscription status: %s", resp.Status)
	}

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

	return resultHeaders, nil
}
