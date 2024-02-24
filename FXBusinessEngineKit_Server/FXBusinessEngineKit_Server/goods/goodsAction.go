package goods

import (
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/tool"
	"FXBusinessEngineKit_Server/user/userController"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Goods struct {
	Receipt string `json:"receipt"`
}

func FetchGoods(content *gin.Context) {

}

func VerifyReceipt(context *gin.Context) {
	transactionParameter, err := tool.BindableModel(VerifyClientParameter{}, context)
	if err != nil {
		tool.ResponseFail(context, tool.HTTPResponseParameterErr, "参数错误", 1)
		return
	}

	userId := tool.UserIdFromHeader(context)

	productModel, productVerifyErr := AppleVerifyReceipt(ProductUrl, transactionParameter, context)
	if productVerifyErr != nil {
		if productModel.Status == 21007 {
			// 切换沙箱验证
			_, sandboxVerifyErr := AppleVerifyReceipt(SandboxUrl, transactionParameter, context)
			if sandboxVerifyErr != nil {
				// 验证失败
				tool.ResponseFail(context, tool.HTTPResponseParameterAppleTransactionVerifyFail, "验证失败", 0)
				return
			}
			// 验证成功
			wallet, sandboxWalletError := userController.FetchUserWallet(userId)
			if sandboxWalletError != nil {
				tool.ResponseFail(context, tool.HTTPResponseParameterAppleTransactionVerifyFail, "订单验证成功，获取用户权益失败", 0)
				return
			}
			walletStr, _ := json.Marshal(wallet)
			tool.ResponseSuccess(context, string(walletStr))
			return
		}
		tool.ResponseFail(context, tool.HTTPResponseParameterAppleTransactionVerifyFail, "验证失败", 1)
		return
	}

	// 验证成功
	wallet, err := userController.FetchUserWallet(userId)
	if err != nil {
		tool.ResponseFail(context, tool.HTTPResponseParameterAppleTransactionVerifyFail, "订单验证成功，获取用户权益失败", 0)
		return
	}
	walletStr, _ := json.Marshal(wallet)
	tool.ResponseSuccess(context, string(walletStr))
}

func VerifyAppleServerNotification(context *gin.Context) {
	bodyData, err := context.GetRawData()
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[苹果交易苹果服务端通知]数据解析失败 %s", err))
		return
	}
	var dataContent map[string]any
	_ = json.Unmarshal(bodyData, &dataContent)

	AppleVerifyNotification(context)
	log.RecordLog(log.Msg, fmt.Sprintf("[苹果交易苹果服务端通知]苹果订单校验信息 %s \n", dataContent))
}
