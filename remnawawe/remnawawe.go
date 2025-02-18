package remnawawe

import (
	"crypto/tls"
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
	client  *http.Client
	baseURL string
	token   string
}

func NewPanel(baseURL, token string) *Panel {
	return &Panel{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		baseURL: baseURL,
		token:   token,
	}
}

func (p *Panel) GetSubscription(shortUuid string) (*SubscriptionResponse, error) {
	httpReq, err := http.NewRequest(http.MethodGet, p.baseURL+"/api/sub/"+shortUuid+"/info", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("error while getting subscription", slog.String("url", p.baseURL+resp.Request.URL.String()))
		return nil, fmt.Errorf("getting subscription status: %s", resp.Status)
	}

	response := &ResponseWrapper{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return &response.Response, nil
}

func (p *Panel) GetUserInfo(shortUuid string) (map[string][]string, error) {
	httpReq, err := http.NewRequest(http.MethodGet, p.baseURL+"/api/sub/"+shortUuid+"/json", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("error while getting subscription", slog.String("url", p.baseURL+resp.Request.URL.String()))
		return nil, fmt.Errorf("getting subscription status: %s", resp.Status)
	}

	filteredHeaders := []string{
		"Profile-Title",
		"Profile-Update-Interval",
		"Subscription-Userinfo",
		"Profile-Web-Page-Url",
	}

	resultHeaders := make(map[string][]string)

	for _, header := range filteredHeaders {
		if values, found := resp.Header[header]; found {
			resultHeaders[header] = values
		}
	}

	return resultHeaders, nil
}
