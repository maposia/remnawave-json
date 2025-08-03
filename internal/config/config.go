package config

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/joho/godotenv"
	"github.com/klauspost/compress/zstd"
)

type config struct {
	remnaweveURL               string
	appHost                    string
	appPort                    string
	webPageTemplate            *template.Template
	happJsonEnabled            bool
	happRouting                string
	httpClient                 *http.Client
	ruOutboundName, ruHostName string
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

func GetHttpClient() *http.Client {
	return conf.httpClient
}

func GetMode() string {
	return os.Getenv("MODE")
}

func GetRuHostName() string {
	return conf.ruHostName
}

func GetRuOutboundName() string {
	return conf.ruOutboundName
}

var conf config

type decompressingRoundTripper struct {
	rt http.RoundTripper
}

func (d *decompressingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Запрашиваем все популярные кодеки
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")

	if GetMode() == "local" {
		req.Header.Set("x-forwarded-for", "127.0.0.1")
		req.Header.Set("x-forwarded-proto", "https")
	}

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

	webPageTemplatePath := os.Getenv("WEB_PAGE_TEMPLATE_PATH")
	if webPageTemplatePath == "" {
		webPageTemplatePath = "/app/templates/subscription/index.html"
	}
	if _, err := os.Stat(webPageTemplatePath); os.IsNotExist(err) {
		slog.Error("File does not exist: " + webPageTemplatePath)
		panic(err)
	}

	conf.happJsonEnabled = os.Getenv("HAPP_JSON_ENABLED") == "true"

	conf.happRouting = os.Getenv("HAPP_ROUTING")

	conf.ruHostName = os.Getenv("RU_USER_HOST")
	conf.ruOutboundName = os.Getenv("RU_OUTBOUND_NAME")

	conf.webPageTemplate, err = template.ParseFiles(webPageTemplatePath)
	if err != nil {
		slog.Error("parsing web page template file:")
		panic(err)
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

func GetRemnawaveToken() any {
	return os.Getenv("REMNAWAVE_TOKEN")
}

func GetMetaTitle() string {
	return os.Getenv("META_TITLE")
}

func GetMetaDescription() string {
	return os.Getenv("META_DESCRIPTION")
}
