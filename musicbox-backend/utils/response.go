package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 结构体定义默认的 JSON 结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Result 统一封装响应方法。
// code 同时作为 HTTP 状态码和响应体 code，保持 REST 语义与前端错误处理一致。
func Result(c *gin.Context, code int, data interface{}, msg string) {
	httpStatus := code
	if http.StatusText(httpStatus) == "" {
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
