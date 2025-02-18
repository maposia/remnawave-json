package utils

import (
	"encoding/json"
	"log"
	"log/slog"
)

func ConvertJsonStringIntoMap(jsonStr string) map[string]interface{} {
	var config map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &config)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}
	return config
}

func DeepCopyMap(original map[string]interface{}) map[string]interface{} {
	data, err := json.Marshal(original)
	if err != nil {
		slog.Error("Error marshaling templates")
		panic(err)
	}

	var copy map[string]interface{}
	err = json.Unmarshal(data, &copy)
	if err != nil {
		slog.Error("Error unmarshaling templates")
		panic(err)
	}

	return copy
}
