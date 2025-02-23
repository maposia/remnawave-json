package config

import (
	"errors"
	"github.com/joho/godotenv"
	"html/template"
	"log/slog"
	"os"
	"remnawave-json/internal/utils"
)

type Config struct {
	V2rayTemplatePath    string
	V2RayTemplate        map[string]interface{}
	V2rayMuxEnabled      bool
	V2rayMuxTemplatePath string
	V2RayMuxTemplate     map[string]interface{}
	RemnaweveURL         string
	APP_HOST             string
	AppPort              string
	WebPageTemplatePath  string
	WebPageTemplate      *template.Template
	HappJsonEnabled      bool
	HappRouting          string
}

var conf Config

func GetConfig() Config {
	return conf
}

func InitConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("Env file not found")
	}

	conf.V2rayTemplatePath = os.Getenv("V2RAY_TEMPLATE_PATH")
	if conf.V2rayTemplatePath == "" {
		conf.V2rayTemplatePath = "/app/templates/v2ray/default.json"
	}
	if _, err := os.Stat(conf.V2rayTemplatePath); os.IsNotExist(err) {
		slog.Error("File does not exist: %s", conf.V2rayTemplatePath)
		panic(err)
	}
	data, err := os.ReadFile(conf.V2rayTemplatePath)
	if err != nil {
		slog.Error("Error reading file:")
		panic(err)
	}
	conf.V2RayTemplate = utils.ConvertJsonStringIntoMap(string(data))

	conf.WebPageTemplatePath = os.Getenv("WEB_PAGE_TEMPLATE_PATH")
	if conf.WebPageTemplatePath == "" {
		conf.WebPageTemplatePath = "/app/templates/subscription/index.html"
	}
	if _, err := os.Stat(conf.WebPageTemplatePath); os.IsNotExist(err) {
		slog.Error("File does not exist: %s", conf.WebPageTemplatePath)
		panic(err)
	}

	conf.HappJsonEnabled = os.Getenv("HAPP_JSON_ENABLED") == "true"

	conf.HappRouting = os.Getenv("HAPP_ROUTING")

	conf.WebPageTemplate, err = template.ParseFiles(conf.WebPageTemplatePath)
	if err != nil {
		slog.Error("parsing web page template file:")
		panic(err)
	}

	conf.V2rayMuxEnabled = os.Getenv("V2RAY_MUX_ENABLED") == "true"

	if conf.V2rayMuxEnabled {
		conf.V2rayMuxTemplatePath = os.Getenv("V2RAY_MUX_TEMPLATE_PATH")
		if conf.V2rayMuxTemplatePath == "" {
			conf.V2rayMuxTemplatePath = "/app/templates/mux/default.json"
		}
		if _, err := os.Stat(conf.V2rayMuxTemplatePath); os.IsNotExist(err) {
			slog.Error("Mux template file does not exist: %s", conf.V2rayMuxTemplatePath)
			panic(err)
		}
		muxData, err := os.ReadFile(conf.V2rayMuxTemplatePath)
		if err != nil {
			slog.Error("Error reading V2ray Mux template file")
			panic(err)
		}
		conf.V2RayMuxTemplate = utils.ConvertJsonStringIntoMap(string(muxData))

	}

	conf.RemnaweveURL = os.Getenv("REMNAWAVE_URL")
	if conf.RemnaweveURL == "" {
		slog.Error("remnawave url not found")
		panic(errors.New("remnawave url not found"))
	}
	conf.APP_HOST = os.Getenv("APP_HOST")
	if conf.APP_HOST == "" {
		conf.APP_HOST = "localhost"
	}
	conf.AppPort = os.Getenv("APP_PORT")
	if conf.AppPort == "" {
		slog.Error("app port not found")
		panic(errors.New("app port not found"))
	}
}
