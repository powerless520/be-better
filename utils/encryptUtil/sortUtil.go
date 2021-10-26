package encryptUtil

import (
	"net/url"
	"sort"
	"strings"
)

type StrPair struct {
	Key string
	Val string
}
type StrPairList []StrPair

func (p StrPairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p StrPairList) Len() int {
	return len(p)
}
func (p StrPairList) Less(i, j int) bool {
	keys := make([]string, 2)
	keys = append(keys, p[i].Val)
	keys = append(keys, p[j].Val)
	sort.Strings(keys)
	if keys[0] == p[0].Val {
		return false
	}
	return true
}

func Ksort(data map[string]string) []string {
	keys := make([]string, 0)
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func KsortWithUrlEncode(params map[string]string) (str string) {
	var key_slice []string
	for k, _ := range params {
		key_slice = append(key_slice, k)
	}
	sort.Strings(key_slice)
	for _, k := range key_slice {
		if k != "sign" {
			str += "&" + k + "=" + url.QueryEscape(params[k])
		}
	}
	str = strings.TrimLeft(str, "&")
	return
}
