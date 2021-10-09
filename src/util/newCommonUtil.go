package util

import (
	"encoding/json"
	"strings"
)

func ParseUriQueryToMap(query string) map[string]string {
	queryMap := strings.Split(query, "&")
	a := make(map[string]string, len(queryMap))
	for _, item := range queryMap {
		itemMap := strings.Split(item, "=")
		a[itemMap[0]] = itemMap[1]
	}
	return a
}

func MapToJson(data map[string]string) string {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonStr)
}
