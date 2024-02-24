package userModel

import (
	"FXBusinessEngineKit_Server/configuration"
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/tool"
	"fmt"
	"time"
)

type FXUserAdapterWalletModel struct {
	TransactionId         string
	OriginalTransactionId string
	ProductId             string
	ExpirationDate        int64
	IsVip                 int
}

type FXUserWallet struct {
	// 用户id
	UserId int64 `json:"userId" gorm:"column:userId"`
	// 设备id
	UUID string `json:"uuid" gorm:"column:uuid"`
	/// 是否是vip
	IsVip int `json:"isVip" gorm:"column:isVIP"`
	/// vip过期时间
	Expires int64 `json:"expires" gorm:"column:expiration"`
	/// 订单信息
	TransactionId string `json:"transactionId" gorm:"column:transactionId"`
	// 订单原始id
	OriginalTransactionId string `json:"originalTransactionId" gorm:"column:originalTransactionId"`
	// 商品id
	GoodsId string `json:"goodsId" gorm:"column:goodsId"`
	// 产品Id
	ProductId int64 `json:"productId" gorm:"column:productId"`
	// 苹果推送消息标识
	AppleNotification int `json:"appleNotification" gorm:"column:appleNotification"`
	// 交易创建时间
	TransactionCreateTime string `json:"transactionCreateTime" gorm:"column:transactionCreateTime"`
	// 每天免费试用次数
	DailyFreeCount int `json:"dailyFreeCount" gorm:"column:dailyFreeCount"`
	// 金币数量
	GoldCoinCount int `json:"goldCoinCount" gorm:"column:goldCoinCount"`
	// 金币充值时间
	GoldCoinTopUpTime string `json:"goldCoinTopUpTime" gorm:"column:goldCoinTopUpTime"`
}

// 如果不加这句代码 GORM 进行操作时表名会加 s

func (u FXUserWallet) TableName() string {
	return configuration.FXDatabaseUserWallet
}

// CheckUserWallet 查询用户钱包
func CheckUserWallet(userId int64) (FXUserWallet, error) {
	var w FXUserWallet
	tool.FXGorm.Where("userId = ?", userId).Find(&w)
	if err := tool.FXGorm.Error; err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[用户信息]用户信息查询权益失败 %s", err))
		return FXUserWallet{}, err
	}
	return w, nil
}

func BandingUserAccountWithTransaction(userId int64, productId int64, transaction FXUserAdapterWalletModel) error {

	date := time.Now().Format(configuration.DateFormat)

	var model FXUserWallet

	tool.FXGorm.Model(&model).Where("userId = ?", userId).Update("isVip", transaction.IsVip)
	tool.FXGorm.Model(&model).Where("userId = ?", userId).Update("expiration", transaction.ExpirationDate)
	tool.FXGorm.Model(&model).Where("userId = ?", userId).Update("transactionId", transaction.TransactionId)
	tool.FXGorm.Model(&model).Where("userId = ?", userId).Update("originalTransactionId", transaction.OriginalTransactionId)
	tool.FXGorm.Model(&model).Where("userId = ?", userId).Update("goodsId", transaction.ProductId)
	tool.FXGorm.Model(&model).Where("userId = ?", userId).Update("transactionCreateTime", date)

	if err := tool.FXGorm.Error; err != nil {
		fmt.Printf("存储异常 %s", err)
		return err
	}

	//w := FXUserWallet{
	//	UserId:                userId,
	//	IsVip:                 transaction.IsVip,
	//	Expires:               transaction.ExpirationDate,
	//	TransactionId:         transaction.TransactionId,
	//	OriginalTransactionId: transaction.OriginalTransactionId,
	//	GoodsId:               transaction.ProductId,
	//	ProductId:             productId,
	//	AppleNotification:     0,
	//	TransactionCreateTime: date,
	//}
	//
	//result := tool.FXGorm.Create(&w)
	//if err := result.Error; err != nil {
	//	log.RecordLog(log.Err, fmt.Sprintf("[创建用户订单]创建用户订单失败 %s", err))
	//	return err
	//}
	//return nil

	return nil
}

func GormUpdateUserWithTransaction(model FXUserWallet) {

	userId := model.UserId

	var m FXUserWallet
	tool.FXGorm.Model(&m).Where("userId = ?", userId).Update("isVip", model.IsVip)
	tool.FXGorm.Model(&m).Where("userId = ?", userId).Update("expiration", model.Expires)
	tool.FXGorm.Model(&m).Where("userId = ?", userId).Update("transactionId", model.TransactionId)
	tool.FXGorm.Model(&m).Where("userId = ?", userId).Update("originalTransactionId", model.OriginalTransactionId)
	tool.FXGorm.Model(&m).Where("userId = ?", userId).Update("goodsId", model.ProductId)
	tool.FXGorm.Model(&m).Where("userId = ?", userId).Update("transactionCreateTime", model.TransactionCreateTime)

	//result := tool.FXGorm.Create(&model)
	//if err := result.Error; err != nil {
	//	log.RecordLog(log.Err, fmt.Sprintf("[苹果服务器通知]更新用户数据，创建订单失败 %s", err))
	//	return
	//}

	log.RecordLog(log.Msg, fmt.Sprintf("[苹果服务器通知]更新用户数据，创建订单成功"))
}
