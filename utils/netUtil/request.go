package netUtil

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/url"
	"strings"
)

func RequestParamsPost(ctx *gin.Context) (params map[string]string) {
	params = map[string]string{}
	for k, v := range ctx.Request.PostForm {
		params[k] = strings.Join(v, "|")
	}
	return params
}

func RequestParamsGet(ctx *gin.Context) (params map[string]string, err error) {
	params = map[string]string{}
	rawQuery := ctx.Request.URL.RawQuery
	if rawQuery == "" {
		return params, errors.New("RequestParamsGet empty")
	}
	m, err := url.ParseQuery(rawQuery)
	if err != nil {
		return params, errors.New("request raw query parse fail")
	}
	for k, v := range m {
		params[k] = v[0]
	}
	return params, nil
}