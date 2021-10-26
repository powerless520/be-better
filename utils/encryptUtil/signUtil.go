package encryptUtil

import (
	"facm/utils/strUtil"
	"strings"
)

func GetSignStr(data map[string]string, includeEmptyParam bool, joinSep string, paramSep string) string {
	sortedKeys := Ksort(data)
	values := make([]string, 0)

	for _, v := range sortedKeys {
		if includeEmptyParam || data[v] != "" {
			values = append(values, v+ paramSep + data[v])
		}
	}

	signStr := strings.Join(values, joinSep)
	return signStr
}

func GetSignStrInterface(data map[string]interface{}, includeEmptyParam bool, joinSep string, paramSep string) string {
	return GetSignStr(strUtil.InterfaceToMap(data), includeEmptyParam, joinSep, paramSep)
}
