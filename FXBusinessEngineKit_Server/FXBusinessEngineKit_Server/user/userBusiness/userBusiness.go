package userBusiness

import (
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/tool"
	"FXBusinessEngineKit_Server/user/userModel"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func UserLoginBusiness(context *gin.Context) (userModel.FXUserLoginModel, error) {
	uuid, err := context.Get("uuid")
	if err == false {
		log.RecordLog(log.Err, fmt.Sprintf("[用户登录]获取uuid失败 %t", err))
		e := errors.New("获取token失败")
		return userModel.FXUserLoginModel{}, e
	}

	token, err := context.Get("token")
	if err == false {
		log.RecordLog(log.Err, fmt.Sprintf("[用户登录]获取token失败 %t", err))
		e := errors.New("获取uuid失败")
		return userModel.FXUserLoginModel{}, e
	}

	productId, err := context.Get("productId")
	if err == false {
		log.RecordLog(log.Err, fmt.Sprintf("[用户登录]获取productId失败 %t", err))
		e := errors.New("获取productId失败")
		return userModel.FXUserLoginModel{}, e
	}

	// 是否需要注册
	u, uErr := registerUserAccountIfNeeded(context)
	if uErr != nil {
		return userModel.FXUserLoginModel{}, uErr
	}

	uuidStr := fmt.Sprintf("%s", uuid)
	tokenStr := fmt.Sprintf("%s", token)
	user := userModel.FXUserLoginModel{
		Uuid:      uuidStr,
		Token:     tokenStr,
		UserId:    u.UserId,
		ProductId: productId.(int64),
		Platform:  u.Platform,
	}

	return user, nil
}

// registerUserAccountIfNeeded 用户登录，库中没有这个用户就自动注册，有这个用户就直接返回 UserId
func registerUserAccountIfNeeded(context *gin.Context) (userModel.FXUser, error) {
	uuid, err := context.Get("uuid")
	if err == false {
		log.RecordLog(log.Err, "[获取用户信息]获取用户信息失败")
		return userModel.FXUser{}, errors.New("获取uuid用户信息失败")
	}

	productId, err := context.Get("productId")
	if err == false {
		log.RecordLog(log.Err, "[获取用户信息]获取productId失败")
		return userModel.FXUser{}, errors.New("获取productId失败")
	}

	platform, err := context.Get("platform")
	if err == false {
		log.RecordLog(log.Err, "[获取用户信息]获取platform失败")
		return userModel.FXUser{}, errors.New("获取platform失败")
	}

	// 类型转换
	uuidStr := fmt.Sprintf("%v", uuid)

	platformStr := fmt.Sprintf("%v", platform)
	// 从数据库获取用户数据
	user, checkErr := userModel.CheckUserIsExist(uuidStr)
	// 新增一个用户
	if checkErr != nil {
		// 查询用户失败，代表是未注册过用户
		model, createErr := userModel.CreateUserAccount(uuidStr, productId.(int64), platformStr)
		if createErr != nil {
			// 失败了，重试一次
			model, _ = userModel.CreateUserAccount(uuidStr, productId.(int64), platformStr)
		}
		jsonStr, err := json.Marshal(model)
		if err != nil {
			// 数据转换失败，返回空
			log.RecordLog(log.Err, fmt.Sprintf("[获取用户信息]获取用户信息成功，数据存储失败 %s", err))
			return userModel.FXUser{}, err
		}
		log.RecordLog(log.Msg, fmt.Sprintf("[获取用户信息]获取用户信息成功 用户信息 %s", jsonStr))
		return model, nil
	}

	return user, nil
}

// FetchUserBusiness 获取用户账户
func FetchUserBusiness(context *gin.Context) (userModel.FXUser, error) {
	userId := tool.UserIdFromHeader(context)

	// 从数据库中找到当前用户
	user, err := userModel.FetchUserModel(userId)
	if err != nil {
		return userModel.FXUser{}, err
	}

	wallet, err := userModel.CheckUserWallet(userId)
	if err != nil {
		return userModel.FXUser{}, err
	}

	user.Wallet = wallet

	return user, nil
}

func FetchUserWalletBusiness(userId int64) (userModel.FXUserWallet, error) {
	w, e := userModel.CheckUserWallet(userId)
	if e != nil {
		return userModel.FXUserWallet{}, e
	}
	return w, nil
}

func UserSaveOrder(userId int64, productId int64, transaction userModel.FXUserAdapterWalletModel) error {
	// 给当前用户创建一条订单信息
	err := userModel.BandingUserAccountWithTransaction(userId, productId, transaction)
	if err != nil {
		return err
	}
	return nil
}
