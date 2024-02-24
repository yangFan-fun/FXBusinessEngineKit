package middleware

import (
	"FXBusinessEngineKit_Server/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func LogMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		// 开始时间
		startTime := time.Now()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		execute := endTime.Sub(startTime)
		// 请求方式
		reqMethod := context.Request.Method
		// 请求路由
		reqUrl := context.Request.RequestURI
		// 状态码
		reqCode := context.Writer.Status()
		// 请求id
		clientIp := context.ClientIP()
		// 请求参数
		reqParams, _ := context.GetRawData()
		// 请求UserAgent
		reqUserAgent := context.Request.UserAgent()
		// 获取uuid
		reqUUID := context.GetHeader("uuid")
		var result logrus.Fields
		result = make(map[string]interface{})
		result["requestMethod"] = reqMethod
		result["requestUrl"] = reqUrl
		result["requestCode"] = reqCode
		result["clientIp"] = clientIp
		result["requestParameters"] = reqParams
		result["userAgent"] = reqUserAgent
		result["execute"] = execute
		result["uuid"] = reqUUID

		log.RecordLog(log.Msg, fmt.Sprintf("[请求记录]%s", result))
	}
}
