package config

import (
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/andybalholm/brotli"
	"github.com/joho/godotenv"
	"github.com/klauspost/compress/zstd"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type config struct {
	v2RayTemplate     map[string]interface{}
	v2rayMuxEnabled   bool
	v2RayMuxTemplate  map[string]interface{}
	remnaweveURL      string
	appHost           string
	appPort           string
	webPageTemplate   *template.Template
	happJsonEnabled   bool
	happRouting       string
	happAnnouncements string
	httpClient        *http.Client
}

func GetHappAnnouncements() string {
	return conf.happAnnouncements
}
func IsHappJsonEnabled() bool {
	return conf.happJsonEnabled
}
func GetHappRouting() string {
	return conf.happRouting
}
func GetAppPort() string {
	return conf.appPort
}
func GetWebPageTemplate() *template.Template {
	return conf.webPageTemplate
}
func GetAppHost() string {
	return conf.appHost
}
func GetRemnaweveURL() string {
	return conf.remnaweveURL
}
func IsMuxEnabled() bool {
	return conf.v2rayMuxEnabled
}
func GetV2RayMuxTemplate() map[string]interface{} {
	return conf.v2RayMuxTemplate
}
func GetV2RayTemplate() map[string]interface{} {
	return conf.v2RayTemplate
}
func GetHttpClient() *http.Client {
	return conf.httpClient
}

var conf config

type decompressingRoundTripper struct {
	rt http.RoundTripper
}

func (d *decompressingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Запрашиваем все популярные кодеки
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")

	resp, err := d.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	encoding := strings.ToLower(resp.Header.Get("Content-Encoding"))
	switch encoding {
	case "gzip":
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body = gr
		resp.Header.Del("Content-Encoding")
	case "deflate":
		resp.Body = flate.NewReader(resp.Body)
		resp.Header.Del("Content-Encoding")
	case "br":
		resp.Body = io.NopCloser(brotli.NewReader(resp.Body))
		resp.Header.Del("Content-Encoding")
	case "zstd":
		dec, err := zstd.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body = io.NopCloser(dec)
		resp.Header.Del("Content-Encoding")
	}

	return resp, nil
}

func InitConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("Env file not found")
	}

	conf.httpClient = &http.Client{
		Transport: &decompressingRoundTripper{
			rt: http.DefaultTransport,
		},
	}

	v2rayTemplatePath := os.Getenv("V2RAY_TEMPLATE_PATH")
	if v2rayTemplatePath == "" {
		v2rayTemplatePath = "/app/templates/v2ray/default.json"
	}
	if _, err := os.Stat(v2rayTemplatePath); os.IsNotExist(err) {
		slog.Error("File does not exist: %s", v2rayTemplatePath)
		panic(err)
	}
	data, err := os.ReadFile(v2rayTemplatePath)
	if err != nil {
		slog.Error("Error reading file:")
		panic(err)
	}
	conf.v2RayTemplate = ConvertJsonStringIntoMap(string(data))

	webPageTemplatePath := os.Getenv("WEB_PAGE_TEMPLATE_PATH")
	if webPageTemplatePath == "" {
		webPageTemplatePath = "/app/templates/subscription/index.html"
	}
	if _, err := os.Stat(webPageTemplatePath); os.IsNotExist(err) {
		slog.Error("File does not exist: %s", webPageTemplatePath)
		panic(err)
	}

	conf.happJsonEnabled = os.Getenv("HAPP_JSON_ENABLED") == "true"

	conf.happRouting = os.Getenv("HAPP_ROUTING")

	announce := os.Getenv("HAPP_ANNOUNCEMENTS")
	conf.happAnnouncements = "base64:" + base64.StdEncoding.EncodeToString([]byte(announce))

	conf.webPageTemplate, err = template.ParseFiles(webPageTemplatePath)
	if err != nil {
		slog.Error("parsing web page template file:")
		panic(err)
	}

	conf.v2rayMuxEnabled = os.Getenv("V2RAY_MUX_ENABLED") == "true"

	if conf.v2rayMuxEnabled {
		v2rayMuxTemplatePath := os.Getenv("V2RAY_MUX_TEMPLATE_PATH")
		if v2rayMuxTemplatePath == "" {
			v2rayMuxTemplatePath = "/app/templates/mux/default.json"
		}
		if _, err := os.Stat(v2rayMuxTemplatePath); os.IsNotExist(err) {
			slog.Error("Mux template file does not exist: %s", v2rayMuxTemplatePath)
			panic(err)
		}
		muxData, err := os.ReadFile(v2rayMuxTemplatePath)
		if err != nil {
			slog.Error("Error reading V2ray Mux template file")
			panic(err)
		}
		conf.v2RayMuxTemplate = ConvertJsonStringIntoMap(string(muxData))

	}

	conf.remnaweveURL = os.Getenv("REMNAWAVE_URL")
	if conf.remnaweveURL == "" {
		slog.Error("remnawave url not found")
		panic(errors.New("remnawave url not found"))
	}
	conf.appHost = os.Getenv("APP_HOST")
	if conf.appHost == "" {
		conf.appHost = "localhost"
	}
	conf.appPort = os.Getenv("APP_PORT")
	if conf.appPort == "" {
		slog.Error("app port not found")
		panic(errors.New("app port not found"))
	}
}

func ConvertJsonStringIntoMap(jsonStr string) map[string]interface{} {
	var config map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &config)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}
	return config
}
