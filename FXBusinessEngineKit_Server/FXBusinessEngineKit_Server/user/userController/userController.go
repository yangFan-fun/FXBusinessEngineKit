package userController

import (
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/tool"
	"FXBusinessEngineKit_Server/user/userBusiness"
	"FXBusinessEngineKit_Server/user/userModel"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserLogin 用户登录
// 返回用户 token uuid userId
func UserLogin(context *gin.Context) {
	user, err := userBusiness.UserLoginBusiness(context)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[用户登录]用户登录失败 %s", err))
		tool.ResponseFail(context, -1, fmt.Sprintf("用户登录失败 %s", err), "")
		return
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[用户登录]用户登录模型序列化失败 %s", err))
		tool.ResponseFail(context, -1, fmt.Sprintf("用户登录模型转换失败 %s", err), "")
		return
	}
	context.JSON(http.StatusOK, gin.H{"code": 200, "msg": "用户注册接口调通", "data": string(userJson)})
}

// User 获取用户信息
func User(context *gin.Context) {
	user, err := userBusiness.FetchUserBusiness(context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"code": 200, "msg": "用户信息获取失败，需要重新登录", "data": nil})
		return
	}
	jsonStr, adapterErr := json.Marshal(user)
	if adapterErr != nil {
		return
	}
	context.JSON(http.StatusOK, gin.H{"code": 200, "msg": "用户信息获取成功", "data": string(jsonStr)})
}

// UserSaveOrder 将数据和用户订单绑定起来
func UserSaveOrder(userId int64, productId int64, transaction userModel.FXUserAdapterWalletModel) error {
	err := userBusiness.UserSaveOrder(userId, productId, transaction)
	if err != nil {
		return err
	}
	return nil
}

// FetchUserWallet 获取用户权益
func FetchUserWallet(userId int64) (userModel.FXUserWallet, error) {
	user, err := userBusiness.FetchUserWalletBusiness(userId)
	if err != nil {
		return user, err
	}
	return user, nil
}
