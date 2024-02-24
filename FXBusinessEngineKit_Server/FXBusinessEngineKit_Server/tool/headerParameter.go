package tool

import (
	"FXBusinessEngineKit_Server/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func TokenFromHeader(context *gin.Context) string {
	token := context.Request.Header.Get("token")
	return token
}

func UUIDFromHeader(context *gin.Context) string {
	uuid := context.Request.Header.Get("uuid")
	return uuid
}

func UserIdFromHeader(context *gin.Context) int64 {
	userId := context.Request.Header.Get("userId")
	ui, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ui = 0000000000
		log.RecordLog(log.Err, fmt.Sprintf("[获取header-userId]字符串转数字失败 %s", err))
	}
	return ui
}

func ProductIdFormHeader(context *gin.Context) int64 {
	productId := context.Request.Header.Get("productId")
	productIdInt, err := strconv.ParseInt(productId, 10, 64)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[获取header-productId]字符串转换失败 %s", err))
		return 000000000
	}
	return productIdInt
}

func PlatformFormHeader(context *gin.Context) string {
	platform := context.Request.Header.Get("platform")
	return platform
}
