package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPResponseCode int64

const (
	HTTPResponseParameterErr                        HTTPResponseCode = 20101
	HTTPResponseParameterDecoderError               HTTPResponseCode = 20102
	HTTPResponseParameterAppleTransactionVerifyFail HTTPResponseCode = 20401
)

func ResponseSuccess(context *gin.Context, data string) {
	context.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": data})
}

func ResponseFail(context *gin.Context, code HTTPResponseCode, msg string, data any) {
	context.JSON(http.StatusBadRequest, gin.H{"code": code, "msg": msg, "data": data})
}
