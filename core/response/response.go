package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JsonResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 1
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, JsonResponse{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func OK(c *gin.Context) {
	Result(SUCCESS, nil, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, nil, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, nil, "fail", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, nil, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, JsonResponse{
		401, nil, "Unauthorized",
	})
}