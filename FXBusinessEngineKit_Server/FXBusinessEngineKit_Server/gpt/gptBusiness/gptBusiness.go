package gptBusiness

import (
	"FXBusinessEngineKit_Server/configuration"
	"FXBusinessEngineKit_Server/gpt/gptModel"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func ChatConnect(context *gin.Context) (gptModel.UserChatModel, error) {

	// 获取参数
	var userModel gptModel.UserChatModel
	err := context.ShouldBind(&userModel)
	if err != nil {
		return gptModel.UserChatModel{}, err
	}

	// 发起请求
	result, err := connect(userModel)
	if err != nil {
		return gptModel.UserChatModel{}, err
	}

	// 获取回答
	answer := result.Choices[0].Message.Content

	// 组装数据
	newUser := gptModel.UserChatModel{Content: answer}
	return newUser, nil
}

func connect(userChat gptModel.UserChatModel) (gptModel.ChatModel, error) {

	requestMsg := []gptModel.ChatRequestModelMessage{
		{
			Role:    "user",
			Content: userChat.Content,
		},
	}

	body := gptModel.ChatRequestModel{
		Messages:         requestMsg,
		Temperature:      0.7,
		TopP:             0.95,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		MaxTokens:        800,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return gptModel.ChatModel{}, err
	}
	payload := strings.NewReader(string(bodyJson))

	request, err := http.NewRequest("POST", configuration.AzureURL, payload)
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("api-key", configuration.AzureKey)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return gptModel.ChatModel{}, err
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		return gptModel.ChatModel{}, err
	}

	var chatModel gptModel.ChatModel
	te := json.Unmarshal(result, &chatModel)
	if te != nil {
		return gptModel.ChatModel{}, te
	}

	defer response.Body.Close()

	return chatModel, nil
}
