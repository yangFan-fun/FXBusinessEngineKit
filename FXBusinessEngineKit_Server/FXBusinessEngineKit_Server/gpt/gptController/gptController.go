package gptController

import (
	"FXBusinessEngineKit_Server/gpt/gptBusiness"
	"FXBusinessEngineKit_Server/tool"
	json2 "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Chat(context *gin.Context) {

	result, err := gptBusiness.ChatConnect(context)
	if err != nil {
		tool.ResponseFail(context, -1, "发起请求失败", fmt.Sprintf("%s", err))
		return
	}

	json, err := json2.Marshal(result)
	if err != nil {
		tool.ResponseFail(context, -1, "模型转换失败", fmt.Sprintf("%s", err))
		return
	}

	tool.ResponseSuccess(context, string(json))
}
