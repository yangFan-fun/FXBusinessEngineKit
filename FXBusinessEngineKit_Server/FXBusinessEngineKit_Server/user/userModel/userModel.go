package userModel

import (
	"FXBusinessEngineKit_Server/configuration"
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/tool"
	"fmt"
	"time"
)

type FXUserLoginModel struct {
	Uuid      string `json:"uuid"`
	Token     string `json:"token"`
	UserId    int64  `json:"userId"`
	ProductId int64  `json:"productId"`
	Platform  string `json:"platform"`
}

type FXUser struct {
	/// 用户id
	UserId int64 `json:"userId" gorm:"column:userId"`
	// 设备id
	UUID string `json:"uuid" gorm:"column:uuid"`
	// 注册时间
	RegistrationDate int64 `json:"registrationDate" gorm:"column:registrationDate"`
	// 产品id
	ProductId int64 `json:"productId" gorm:"column:productId"`
	// 平台
	Platform string `json:"platform" gorm:"column:platform"`
	/// 用户钱包信息
	Wallet FXUserWallet `json:"wallet" gorm:"-"`
}

// 如果不加这句代码 GORM 进行操作时表名会加 s

func (u FXUser) TableName() string {
	return configuration.FXDatabaseUser
}

func CreateUserAccount(uuid string, productId int64, platform string) (FXUser, error) {
	// 定义一个用户并初始化数据

	w := FXUserWallet{}

	u := FXUser{
		UUID:             uuid,
		Wallet:           w,
		ProductId:        productId,
		RegistrationDate: time.Now().Unix(),
		Platform:         platform,
	}

	// 插入一条数据
	// 代码会自动生成 SQL语句 insert into `fxuser` (`key`) valuers (value)
	result := tool.FXGorm.Create(&u)
	if err := result.Error; err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[用户信息]新用户数据插入数据失败%s", err))
		return FXUser{}, err
	}

	// 从数据库中获取用户数据
	nu, err := FindUser(uuid)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[用户信息]新用户创建查询用户数据失败 %s", err))
		return FXUser{}, err
	}

	log.RecordLog(log.Msg, fmt.Sprintf("[用户信息]新用户数据库插入数据成功"))

	// 获取当前用户
	current, err := CheckUserIsExist(uuid)
	currentUserWallet := FXUserWallet{
		UserId:                current.UserId,
		UUID:                  current.UUID,
		IsVip:                 0,
		Expires:               0,
		TransactionId:         "",
		OriginalTransactionId: "",
		GoodsId:               "",
		ProductId:             current.ProductId,
		AppleNotification:     0,
		TransactionCreateTime: "",
		DailyFreeCount:        2,
		GoldCoinCount:         0,
		GoldCoinTopUpTime:     "",
	}
	wr := tool.FXGorm.Create(&currentUserWallet)
	if err := wr.Error; err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[用户信息]新用户插入金额数据失败 %s", err))
	}

	return nu, nil
}

// CheckUserIsExist 检查用户是否已经注册
// 已经注册返回用户信息，未注册返回 false
func CheckUserIsExist(uuid string) (FXUser, error) {
	u := FXUser{}
	// 获取第一个匹配记录
	// select * from fxuser where uuid = uuid limit 1
	result := tool.FXGorm.Where("uuid = ?", uuid).First(&u)
	if err := result.Error; err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[用户信息]查询用户是否存在，用户不存在 %s", err))
		return u, err
	}

	log.RecordLog(log.Msg, fmt.Sprintf("[用户信息]查询用户是否存在，查询成功 uuid:%s", u.UUID))
	return u, nil
}

// FindUser 找到指定用户
func FindUser(uuid string) (FXUser, error) {
	log.RecordLog(log.Msg, fmt.Sprintf("[数据库]正在准备查询指定用户"))
	u := FXUser{}
	result := tool.FXGorm.Where("uuid = ?", uuid).First(&u)
	if err := result.Error; err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[数据库]数据库查询数据失败%s", err))
		return u, err
	}
	return u, nil
}

// FetchUserModel 从数据库中获取用户
func FetchUserModel(userId int64) (FXUser, error) {
	u := FXUser{}

	tool.FXGorm.Where("userId = ?", userId).First(&u)
	if err := tool.FXGorm.Error; err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[数据库]查询用户失败 %s", err))
		return FXUser{}, err
	}

	// 查询用户权益
	w, err := CheckUserWallet(userId)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[数据库]查询用户权益失败 %s", err))
		return FXUser{}, err
	}

	u.Wallet = w

	userSubscribeStatusUpdateIfNeeded(u)

	return u, nil
}

func userSubscribeStatusUpdateIfNeeded(user FXUser) {
	// 如果用户是会员 检查用户的过期时间
	if user.Wallet.IsVip == 1 {
		currentDate := time.Now().Unix()
		// 会员过期
		if currentDate >= user.Wallet.Expires {
			GormUpdateUserVipStatus(user)
		}
	}
}

// GormUpdateUserVipStatus 设置用户会员状态
func GormUpdateUserVipStatus(user FXUser) {
	// 更新指定字段的数据
	// update `fxuser` set `isVip` = isVip where uuid = uuid

	// 更新会员状态
	tool.FXGorm.Model(&FXUserWallet{}).Where("uuid = ?", user.UUID).Update("isVip", 1)
	tool.FXGorm.Model(&FXUserWallet{}).Where("uuid = ?", user.UUID).Update("vipExpiration", user.Wallet.Expires)
	log.RecordLog(log.Msg, fmt.Sprintf("[数据库]用户会员信息更新，用户是会员，过期时间%d", user.Wallet.Expires))
}
