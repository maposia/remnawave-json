package service

import (
	"encoding/base64"
	"encoding/json"
	"github.com/xtls/libxray"
	"log/slog"
	"net/http"
	"remnawave-json/internal/config"
	"remnawave-json/internal/utils"
	"remnawave-json/remnawave"
)

type Service struct {
	Panel remnawave.Panel
}

func (s *Service) GenerateJson(shortUuid string) ([]interface{}, http.Header, error) {
	sub, err := s.Panel.GetSubscription(shortUuid)
	headers, _ := s.Panel.GetUserInfo(shortUuid)

	if err != nil {
		slog.Error("Get Subscription Error", err)
		return nil, nil, err
	}

	var jsonSub []interface{}

	for _, link := range sub.Links {
		encodedLink := base64.StdEncoding.EncodeToString([]byte(link))
		encodedJson, err := base64.StdEncoding.DecodeString(libXray.ConvertShareLinksToXrayJson(encodedLink))
		if err != nil {
			slog.Error("error while decoding base64 link")
			panic(err)
		}
		jsonConf := utils.ConvertJsonStringIntoMap(string(encodedJson))["data"].(map[string]interface{})

		configCopy := utils.DeepCopyMap(config.GetConfig().V2RayTemplate)

		newOutbounds := jsonConf["outbounds"]

		if outboundsArray, ok := newOutbounds.([]interface{}); ok {
			for _, outbound := range outboundsArray {
				if outboundMap, ok := outbound.(map[string]interface{}); ok {
					if sendThrough, exists := outboundMap["sendThrough"]; exists {
						configCopy["remarks"] = convertStringToUnicodeEscaped(sendThrough.(string))
						delete(outboundMap, "sendThrough")
					}
				}
			}
		}

		if outboundsArray, ok := newOutbounds.([]interface{}); ok {
			for _, outbound := range outboundsArray {
				if outboundMap, ok := outbound.(map[string]interface{}); ok {
					outboundMap["tag"] = "proxy"
					outboundMap["mux"] = config.GetConfig().V2RayMuxTemplate
				}
			}
		}

		if outbounds, ok := configCopy["outbounds"].([]interface{}); ok {
			for _, newOutbound := range newOutbounds.([]interface{}) {
				outbounds = append([]interface{}{newOutbound}, outbounds...) // Prepend to the beginning
			}
			configCopy["outbounds"] = outbounds
		} else {
			configCopy["outbounds"] = newOutbounds
		}

		cleanJsonData(configCopy)
		jsonSub = append(jsonSub, cleanJsonData(configCopy))
	}

	return jsonSub, headers, nil
}

func convertStringToUnicodeEscaped(input string) string {
	encoded, err := json.Marshal(input)
	if err != nil {
		return "proxy"
	}
	return string(encoded[1 : len(encoded)-1])
}

func cleanJsonData(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = cleanJsonData(value)
			if v[key] == nil || v[key] == "" || isZero(v[key]) || isEmptyMap(v[key]) || isEmptySlice(v[key]) {
				delete(v, key)
			}
		}
		return v

	case []interface{}:
		var newArray []interface{}
		for _, item := range v {
			cleanedItem := cleanJsonData(item)
			if cleanedItem != nil && cleanedItem != "" && !isZero(cleanedItem) && !isEmptyMap(cleanedItem) && !isEmptySlice(cleanedItem) {
				newArray = append(newArray, cleanedItem)
			}
		}
		return newArray

	default:
		return v
	}
}

func isZero(value interface{}) bool {
	switch v := value.(type) {
	case int:
		return v == 0
	case int64:
		return v == 0
	case float64:
		return v == 0.0
	default:
		return false
	}
}

func isEmptyMap(data interface{}) bool {
	if m, ok := data.(map[string]interface{}); ok {
		return len(m) == 0
	}
	return false
}

func isEmptySlice(data interface{}) bool {
	if s, ok := data.([]interface{}); ok {
		return len(s) == 0
	}
	return false
}

func NewService(panel *remnawave.Panel) *Service {
	return &Service{
		Panel: *panel,
	}
}
