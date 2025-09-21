package remnawave

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"remnawave-json/internal/config"
)

type ResponseConverterWrapper struct {
	Response XrayConverterResponse `json:"response"`
}

type XrayConverterResponse struct {
	User              UserRaw           `json:"user"`
	ConvertedUserInfo ConvertedUserInfo `json:"convertedUserInfo"`
	Headers           map[string]string `json:"headers"`
	RawHosts          []RawHost         `json:"rawHosts"`
}

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
	Address           string            `json:"address"`
	ALPN              string            `json:"alpn"`
	Fingerprint       string            `json:"fingerprint"`
	Host              string            `json:"host"`
	Network           string            `json:"network"`
	Password          Passwords         `json:"password"`
	Path              string            `json:"path"`
	PublicKey         string            `json:"publicKey"`
	Port              int               `json:"port"`
	Protocol          string            `json:"protocol"`
	Remark            string            `json:"remark"`
	ShortID           string            `json:"shortId"`
	SNI               string            `json:"sni"`
	SpiderX           string            `json:"spiderX"`
	TLS               string            `json:"tls"`
	HeaderType        *string           `json:"headerType"`
	AdditionalParams  *AdditionalParams `json:"additionalParams"`
	XHttpExtraParams  *interface{}      `json:"xHttpExtraParams"`
	MuxParams         *interface{}      `json:"muxParams"`
	SockoptParams     *interface{}      `json:"sockoptParams"`
	ServerDescription string            `json:"serverDescription"`
	Flow              string            `json:"flow"`
	AllowInsecure     bool              `json:"allowInsecure"`
	ProtocolOptions   *ProtocolOptions  `json:"protocolOptions"`
	DbData            DbData            `json:"dbData"`
}

type AdditionalParams struct {
	Mode            *string  `json:"mode"`
	HeartbeatPeriod *float64 `json:"heartbeatPeriod"`
}

type ProtocolOptions struct {
	SS *SSOptions `json:"ss"`
}

type SSOptions struct {
	Method *string `json:"method"`
}

type DbData struct {
	RawInbound               *interface{} `json:"rawInbound"`
	InboundTag               string       `json:"inboundTag"`
	UUID                     string       `json:"uuid"`
	ConfigProfileUUID        *string      `json:"configProfileUuid"`
	ConfigProfileInboundUUID *string      `json:"configProfileInboundUuid"`
	IsDisabled               bool         `json:"isDisabled"`
	ViewPosition             float64      `json:"viewPosition"`
	Remark                   string       `json:"remark"`
	IsHidden                 bool         `json:"isHidden"`
	Tag                      *string      `json:"tag"`
	VlessRouteID             *int         `json:"vlessRouteId"`
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

type XrayConfig struct {
	DNS       DNS        `json:"dns"`
	Inbounds  []Inbound  `json:"inbounds"`
	Outbounds []Outbound `json:"outbounds"`
	Remarks   string     `json:"remarks"`
	Routing   Routing    `json:"routing"`
}

type DNS struct {
	QueryStrategy string        `json:"queryStrategy"`
	Servers       []interface{} `json:"servers"`
}

type FreedomSettings struct{}

type BlackholeSettings struct{}

type Inbound struct {
	Listen   string      `json:"listen"`
	Port     int         `json:"port"`
	Protocol string      `json:"protocol"`
	Settings interface{} `json:"settings"`
	Sniffing Sniffing    `json:"sniffing"`
	Tag      string      `json:"tag"`
}

type Outbound struct {
	Protocol       string          `json:"protocol"`
	Settings       interface{}     `json:"settings,omitempty"`
	StreamSettings *StreamSettings `json:"streamSettings,omitempty"`
	Tag            string          `json:"tag"`
}

type StreamSettings struct {
	Network         string          `json:"network"`
	Security        string          `json:"security"`
	RealitySettings RealitySettings `json:"realitySettings"`
	TCPSettings     interface{}     `json:"tcpSettings"`
}

type RealitySettings struct {
	Fingerprint string `json:"fingerprint"`
	PublicKey   string `json:"publicKey"`
	ServerName  string `json:"serverName"`
	ShortID     string `json:"shortId"`
	Show        bool   `json:"show"`
}

type VLESSSettings struct {
	VNext []VNext `json:"vnext"`
}

type VNext struct {
	Address string      `json:"address"`
	Port    int         `json:"port"`
	Users   []UserToRaw `json:"users"`
}
type Routing struct {
	DomainMatcher  string      `json:"domainMatcher"`
	DomainStrategy string      `json:"domainStrategy"`
	Rules          []Rule      `json:"rules"`
	Balancers      []Balancer  `json:"balancers"`
	Observatory    Observatory `json:"observatory"`
}

type Rule struct {
	Type        string   `json:"type"`
	InboundTag  []string `json:"inboundTag,omitempty"`
	OutboundTag string   `json:"outboundTag,omitempty"`
	BalancerTag string   `json:"balancerTag,omitempty"`
	Protocol    []string `json:"protocol,omitempty"`
	Domain      []string `json:"domain,omitempty"`
	IP          []string `json:"ip,omitempty"`
	Port        string   `json:"port,omitempty"`
}

type Balancer struct {
	Tag      string   `json:"tag"`
	Selector []string `json:"selector"`
	Strategy Strategy `json:"strategy"`
}

type UserRaw struct {
	UUID                     string                `json:"uuid"`
	ShortUUID                string                `json:"shortUuid"`
	Username                 string                `json:"username"`
	Status                   string                `json:"status"`
	UsedTrafficBytes         float64               `json:"usedTrafficBytes"`
	LifetimeUsedTrafficBytes float64               `json:"lifetimeUsedTrafficBytes"`
	TrafficLimitBytes        int64                 `json:"trafficLimitBytes"`
	TrafficLimitStrategy     string                `json:"trafficLimitStrategy"`
	SubLastUserAgent         *string               `json:"subLastUserAgent"`
	SubLastOpenedAt          *string               `json:"subLastOpenedAt"`
	ExpireAt                 string                `json:"expireAt"`
	OnlineAt                 *string               `json:"onlineAt"`
	SubRevokedAt             *string               `json:"subRevokedAt"`
	LastTrafficResetAt       *string               `json:"lastTrafficResetAt"`
	TrojanPassword           string                `json:"trojanPassword"`
	VlessUUID                string                `json:"vlessUuid"`
	SSPassword               string                `json:"ssPassword"`
	Description              *string               `json:"description"`
	Tag                      *string               `json:"tag"`
	TelegramID               *int64                `json:"telegramId"`
	Email                    *string               `json:"email"`
	HwidDeviceLimit          *int                  `json:"hwidDeviceLimit"`
	FirstConnectedAt         *string               `json:"firstConnectedAt"`
	LastTriggeredThreshold   int                   `json:"lastTriggeredThreshold"`
	CreatedAt                string                `json:"createdAt"`
	UpdatedAt                string                `json:"updatedAt"`
	ActiveInternalSquads     []ActiveInternalSquad `json:"activeInternalSquads"`
	SubscriptionURL          string                `json:"subscriptionUrl"`
	LastConnectedNode        *LastConnectedNode    `json:"lastConnectedNode"`
	Happ                     Happ
}
type ConvertedUserInfo struct {
	DaysLeft            float64 `json:"daysLeft"`
	TrafficLimit        string  `json:"trafficLimit"`
	TrafficUsed         string  `json:"trafficUsed"`
	LifetimeTrafficUsed string  `json:"lifetimeTrafficUsed"`
	IsHwidLimited       bool    `json:"isHwidLimited"`
}

type ActiveInternalSquad struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type LastConnectedNode struct {
	ConnectedAt string `json:"connectedAt"`
	NodeName    string `json:"nodeName"`
	CountryCode string `json:"countryCode"`
}

type Happ struct {
	CryptoLink string `json:"cryptoLink"`
}

type Strategy struct {
	Type string `json:"type"`
}

type Observatory struct {
	SubjectSelector   []string `json:"subjectSelector"`
	ProbeUrl          string   `json:"probeUrl"`
	ProbeInterval     string   `json:"probeInterval"`
	EnableConcurrency bool     `json:"enableConcurrency"`
}
type Sniffing struct {
	DestOverride []string `json:"destOverride"`
	Enabled      bool     `json:"enabled"`
	RouteOnly    bool     `json:"routeOnly,omitempty"`
}

type UserToRaw struct {
	Encryption string `json:"encryption"`
	Flow       string `json:"flow"`
	ID         string `json:"id"`
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

func GetRawSubscription(shortUuid string, header string) (*ResponseConverterWrapper, error) {
	url := fmt.Sprintf("%s/api/subscriptions/by-short-uuid/%s/raw", config.GetRemnaweveURL(), shortUuid)
	slog.Info("Making request to raw subscription", "url", url, "shortUuid", shortUuid)

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", header)

	token := config.GetRemnawaveToken()
	if token != nil && token != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := config.GetHttpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("error while getting subscription", slog.String("url", config.GetRemnaweveURL()+resp.Request.URL.String()))
		return nil, fmt.Errorf("getting subscription status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var wrapper ResponseConverterWrapper
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &wrapper, nil
}

func ConvertToXrayConfig(wrapper *ResponseConverterWrapper) ([]byte, error) {
	response := wrapper.Response

	if len(response.RawHosts) == 0 {
		return nil, fmt.Errorf("no raw hosts found in response")
	}
	var remarks string
	if response.RawHosts[0].Remark != "" {
		remarks = response.RawHosts[0].Remark
	} else {
		remarks = "Сервер VPN"
	}

	config := XrayConfig{
		DNS: DNS{
			QueryStrategy: "UseIPv4",
			Servers: []interface{}{
				"94.140.14.14",
				map[string]interface{}{
					"address": "94.140.14.14",
					"domains": []string{"geosite:youtube", "geosite:category-ban-ru"},
					"port":    53,
				},
				map[string]interface{}{
					"address": "94.140.15.15",
					"domains": []string{"geosite:private", "geosite:category-ru", "geosite:apple", "geosite:twitch"},
					"port":    53,
				},
			},
		},
		Inbounds: []Inbound{
			{
				Listen:   "127.0.0.1",
				Port:     10808,
				Protocol: "socks",
				Settings: map[string]interface{}{
					"auth": "noauth",
					"udp":  true,
				},
				Sniffing: Sniffing{
					DestOverride: []string{"http", "tls", "quic"},
					Enabled:      true,
					RouteOnly:    false,
				},
				Tag: "socks",
			},
			{
				Listen:   "127.0.0.1",
				Port:     10809,
				Protocol: "http",
				Settings: map[string]interface{}{
					"allowTransparent": false,
				},
				Sniffing: Sniffing{
					DestOverride: []string{"http", "tls", "quic"},
					Enabled:      true,
					RouteOnly:    false,
				},
				Tag: "http",
			},
		},
		Remarks: remarks,
		Routing: Routing{
			DomainMatcher:  "hybrid",
			DomainStrategy: "IPIfNonMatch",
			Rules: []Rule{
				{
					Type: "field",
					Domain: []string{
						"geosite:private",
						"geosite:category-ru",
						"geosite:apple",
						"geosite:twitch",
					},
					OutboundTag: "direct",
				},
				{
					Type: "field",
					IP: []string{
						"geoip:ru",
						"geoip:private",
					},
					OutboundTag: "direct",
				},
				{
					Type: "field",
					Domain: []string{
						"geosite:youtube",
						"geosite:category-ban-ru",
					},
					BalancerTag: "proxy-balancer",
				},
				{
					Type: "field",
					IP: []string{
						"158.85.224.160/27",
						"158.85.46.128/27",
						"158.85.5.192/27",
						"173.192.222.160/27",
						"173.192.231.32/27",
						"18.194.0.0/15",
						"184.173.128.0/17",
						"208.43.122.128/27",
						"34.224.0.0/12",
						"50.22.198.204/30",
						"54.242.0.0/15",
						"91.108.56.0/22",
						"91.108.4.0/22",
						"91.108.8.0/22",
						"91.108.16.0/22",
						"91.108.12.0/22",
						"149.154.160.0/20",
						"91.105.192.0/23",
						"91.108.20.0/22",
						"85.76.151.0/24",
						"2001:b28:f23d::/48",
						"2001:b28:f23f::/48",
						"2001:67c:4e8::/48",
						"2001:b28:f23c::/48",
						"2a0a:f280::/32",
					},
					BalancerTag: "proxy-balancer",
				},
				{
					Type: "field",
					IP: []string{
						"94.140.14.14",
					},
					Port:        "53",
					BalancerTag: "proxy-balancer",
				},
				{
					Type: "field",
					IP: []string{
						"94.140.15.15",
					},
					Port:        "53",
					OutboundTag: "direct",
				},
				// Правило для socks-direct inbound
				{
					Type: "field",
					InboundTag: []string{
						"socks-direct",
					},
					OutboundTag: "direct",
				},
				{
					Type:        "field",
					InboundTag:  []string{"socks", "http"},
					BalancerTag: "proxy-balancer",
				},
			},
			Balancers: []Balancer{
				{
					Tag:      "proxy-balancer",
					Selector: []string{"proxy1", "proxy2"},
					Strategy: Strategy{Type: "roundRobin"},
				},
			},
			Observatory: Observatory{
				SubjectSelector:   []string{"proxy"},
				ProbeUrl:          "https://www.google.com/generate_204",
				ProbeInterval:     "10s",
				EnableConcurrency: false,
			},
		},
	}

	hostsToProcess := len(response.RawHosts)
	if hostsToProcess > 2 {
		hostsToProcess = 2
	}

	for i := 0; i < hostsToProcess; i++ {
		host := response.RawHosts[i]

		if host.Protocol != "vless" {
			return nil, fmt.Errorf("host %d must use vless protocol", i+1)
		}
		if host.TLS != "reality" {
			return nil, fmt.Errorf("host %d must use reality TLS", i+1)
		}
		if host.Network != "tcp" {
			return nil, fmt.Errorf("host %d must use tcp network", i+1)
		}

		fingerprint := "chrome"
		if host.Fingerprint != "" {
			fingerprint = host.Fingerprint
		}

		outbound := Outbound{
			Protocol: "vless",
			Settings: VLESSSettings{
				VNext: []VNext{
					{
						Address: host.Address,
						Port:    host.Port,
						Users: []UserToRaw{
							{
								Encryption: "none",
								Flow:       host.Flow,
								ID:         host.Password.VlessPassword,
							},
						},
					},
				},
			},
			StreamSettings: &StreamSettings{
				Network:  "tcp",
				Security: "reality",
				RealitySettings: RealitySettings{
					Fingerprint: fingerprint,
					PublicKey:   host.PublicKey,
					ServerName:  host.SNI,
					ShortID:     host.ShortID,
					Show:        false,
				},
				TCPSettings: map[string]interface{}{},
			},
			Tag: fmt.Sprintf("proxy%d", i+1),
		}
		config.Outbounds = append(config.Outbounds, outbound)
	}

	config.Outbounds = append(config.Outbounds,
		Outbound{Protocol: "freedom", Tag: "direct", Settings: FreedomSettings{}},
		Outbound{Protocol: "blackhole", Tag: "block", Settings: BlackholeSettings{}},
		Outbound{Protocol: "blackhole", Tag: "TORRENT", Settings: BlackholeSettings{}},
	)

	outputJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal output JSON: %w", err)
	}

	return outputJSON, nil
}
