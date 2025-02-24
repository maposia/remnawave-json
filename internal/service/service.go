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

func (s *Service) GenerateJson(shortUuid string, header string) ([]interface{}, http.Header, error) {
	headers, body, _ := s.Panel.GetUserInfo(shortUuid, header)

	encJson, err := base64.StdEncoding.DecodeString(libXray.ConvertShareLinksToXrayJson(body))
	if err != nil {
		slog.Error("decode xray json config error", err)
		return nil, nil, err
	}
	outbounds := utils.ConvertJsonStringIntoMap(string(encJson))["data"].(map[string]interface{})["outbounds"].([]interface{})

	jsonSub := make([]interface{}, len(outbounds))

	for i, outbound := range outbounds {
		configCopy := utils.DeepCopyMap(config.GetConfig().V2RayTemplate)

		if outboundMap, ok := outbound.(map[string]interface{}); ok {
			outboundMap["tag"] = "proxy"

			if sendThrough, exists := outboundMap["sendThrough"]; exists {
				configCopy["remarks"] = convertStringToUnicodeEscaped(sendThrough.(string))
				delete(outboundMap, "sendThrough")
			}

			if config.GetConfig().V2rayMuxEnabled {
				outboundMap["mux"] = config.GetConfig().V2RayMuxTemplate
			}
		}

		if existingOutbounds, ok := configCopy["outbounds"].([]interface{}); ok {
			configCopy["outbounds"] = append([]interface{}{outbound}, existingOutbounds...)
		} else {
			configCopy["outbounds"] = []interface{}{outbound}
		}

		jsonSub[i] = cleanJsonData(configCopy)
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
