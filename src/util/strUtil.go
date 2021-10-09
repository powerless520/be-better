package util

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func FilterAlphaNumberChinese(str string, extra string) string {
	if str == "" {
		return str
	}
	str_new := strings.FieldsFunc(str, func(ch rune) bool {
		switch {
		case ch >= '0' && ch <= '9':
			return false
		case ch >= 'a' && ch <= 'z':
			return false
		case ch >= 'A' && ch <= 'Z':
			return false
		case unicode.Is(unicode.Scripts["Han"], ch):
			return false
		case extra != "" && strings.Contains(extra, string(ch)):
			return false
		}
		return true
	})
	return strings.Join(str_new, "")
}
func FilterAlphaNumber(str string, extra string) string {
	if str == "" {
		return str
	}
	str_new := strings.FieldsFunc(str, func(ch rune) bool {
		switch {
		case ch >= '0' && ch <= '9':
			return false
		case ch >= 'a' && ch <= 'z':
			return false
		case ch >= 'A' && ch <= 'Z':
			return false
		case extra != "" && strings.Contains(extra, string(ch)):
			return false
		}
		return true
	})
	return strings.Join(str_new, "")
}

func Page_Html(total_cnt, page_rows, current_page int) template.HTML {
	var page_html string = ""
	if current_page < 1 {
		current_page = 1
	}

	total_cnt64 := float64(total_cnt)
	page_rows64 := float64(page_rows)
	total_pages64 := total_cnt64 / page_rows64
	total_pages := int(math.Ceil(total_pages64))

	start_page := current_page - 3
	end_page := current_page + 3

	if total_pages <= 7 {
		start_page = 1
		end_page = total_pages
	} else {
		if current_page-3 <= 0 {
			start_page = 1
			end_page = 7
		}
		if total_pages-current_page < 3 {
			end_page = total_pages
			start_page = total_pages - 7
		}
	}

	page_html += "<ul>"
	if current_page <= 1 {
		page_html += "<li class=\"disabled\"><a href=\"#\">首页</a></li>"
	} else {
		page_html += "<li><a href=\"#\" onclick=\"return page_jump(1);\">首页</a></li>"
	}

	var active string
	for i := start_page; i <= end_page; i++ {
		active = ""
		if i == current_page {
			active = "active"
		}
		page_html += "<li class=\"" + active + "\"><a href=\"#\" onclick=\"return page_jump(" + strconv.Itoa(i) + ");\">" + strconv.Itoa(i) + "</a></li>"
	}

	if current_page >= total_pages {
		page_html += "<li class=\"disabled\"><a href=\"#\">末页</a></li>"
	} else {
		page_html += "<li><a href=\"#\" onclick=\"return page_jump(" + strconv.Itoa(total_pages) + ");\">末页</a></li>"
	}
	page_html += "</ul>"

	return template.HTML(page_html)
}

func HttpBuildQuery(data map[string]string) string {
	sorted_keys := make([]string, 0)
	for k, _ := range data {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string
	for i, k := range sorted_keys {
		value := fmt.Sprintf("%v", data[k])
		if i != (len(sorted_keys) - 1) {
			signStrings = signStrings + k + "=" + value + "&"
		} else {
			signStrings = signStrings + k + "=" + value
		}
	}
	return signStrings
}
func GetRandomStr(l int) string {
	str_byte := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, str_byte[r.Intn(len(str_byte))])
	}
	return string(result)
}

func FilterDelHtmlLabel(htmlstr string, length_jump int, del_chars string) string {
	if htmlstr == "" {
		return ""
	}
	if length_jump > 0 && len(htmlstr) < length_jump {
		return htmlstr
	}
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	htmlstr = re.ReplaceAllStringFunc(htmlstr, strings.ToLower)
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	htmlstr = re.ReplaceAllString(htmlstr, "")
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	htmlstr = re.ReplaceAllString(htmlstr, "")
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	htmlstr = re.ReplaceAllString(htmlstr, "\n")
	re, _ = regexp.Compile("\\s{2,}")
	htmlstr = re.ReplaceAllString(htmlstr, "\n")
	if del_chars != "" {
		for _, ch := range del_chars {
			htmlstr = strings.Replace(htmlstr, string(ch), "", -1)
		}
	}
	return strings.TrimSpace(htmlstr)
}

func MergeJson(dst, src map[string]interface{}, ignoreEmpty bool) map[string]interface{} {
	return innerMergeJson(dst, src, 0, ignoreEmpty)
}

func innerMergeJson(dst, src map[string]interface{}, depth int, ignoreEmpty bool) map[string]interface{} {
	if dst == nil {
		dst = make(map[string]interface{})
	}
	if depth > 32 {
		return dst
		// panic("too deep!")
	}

	for key, srcVal := range src {

		if dstVal, ok := dst[key]; ok {

			srcMap, srcMapOk := innerJsonMapify(srcVal, ignoreEmpty)
			dstMap, dstMapOk := innerJsonMapify(dstVal, ignoreEmpty)

			if srcMapOk && dstMapOk {
				srcVal = innerMergeJson(dstMap, srcMap, depth+1, ignoreEmpty)
			}
		}

		if !ignoreEmpty || srcVal != nil {
			value := reflect.ValueOf(srcVal)
			if value.Kind() == reflect.String {
				if srcVal != "" {
					dst[key] = srcVal
				}
			} else {
				dst[key] = srcVal
			}
		}
	}

	return dst
}

func innerJsonMapify(i interface{}, ignoreEmpty bool) (map[string]interface{}, bool) {

	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Map {

		m := map[string]interface{}{}

		for _, k := range value.MapKeys() {
			if !ignoreEmpty || value.MapIndex(k).Interface() != nil {
				m[k.String()] = value.MapIndex(k).Interface()
			}
		}

		return m, true
	}

	return map[string]interface{}{}, false
}

// CompareMoney 校验金额
func CompareMoney(money1, money2 string) error {

	floatMoney1, err := strconv.ParseFloat(money1, 64)
	if err != nil {
		return errors.New("parseFloat1Err:" + money1 + "|" + err.Error())
	}

	floatMoney2, err := strconv.ParseFloat(money2, 64)
	if err != nil {
		return errors.New("parseFloat2Err:" + money2 + "|" + err.Error())
	}

	if floatMoney1 != floatMoney2 {
		return errors.New("compareMoneyErr:" + money1 + "|" + money2)
	}

	return nil
}
