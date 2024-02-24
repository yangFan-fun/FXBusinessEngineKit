package configuration

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

//---------------------------------------------
// FXBusinessEngineConfiguration 应用配置相关内容
//---------------------------------------------

// AppPort 端口
var AppPort string = "8080"

// DateFormat 日期格式
var DateFormat = "2006-01-02 15:04:05"

// FXDatabaseServer 数据库连接信息
var FXDatabaseServer = ""

// FXDatabaseUser 用户信息数据库名称
var FXDatabaseUser = "FXBusinessEngineDatabase_user"

// FXDatabaseUserWallet 用户钱包数据库名称
var FXDatabaseUserWallet = "FXBusinessEngineDatabase_wallet"

// FXDatabaseProduct 应用产品标识数据库名称
var FXDatabaseProduct = "FXBusinessEngineDatabase_product"

//---------------------------------------------
// JWT (Json Web Token)
// FXBusinessEngineKit_Server 接口鉴权采用JWT方案
//---------------------------------------------

// JWTExp 过期时间30天
var JWTExp = time.Hour * 24 * 30

// JWTSecret 密钥
var JWTSecret = "FXBusinessEngineProject"

type JWTHeader struct {
	typ string
	alg string
}

type MyClaims struct {
	UserId string `json: "userId"`
	jwt.RegisteredClaims
}

//--------------------------------------------------------------------
// SSL证书配置内容
// SSL证书与HTTPS的关系：
// 基于 SSL 证书，可将站点由 HTTP（Hypertext Transfer Protocol）
// 切换到 HTTPS（Hyper Text Transfer Protocol over Secure Socket Layer）
// 即基于安全套接字层（SSL）进行安全数据传输的加密版 HTTP 协议
//--------------------------------------------------------------------

// HttpsCertFile HTTPS证书文件名
var HttpsCertFile = ""

// HttpsKeyFile Https证书文件名
var HttpsKeyFile = ""

//----------------------------------------------------------
// 腾讯云对象存储COS配置内容
// FXBusinessEngineKit_Server已接入腾讯云对象存储COS的GO语言SDK
// 只需要在腾讯云官网上开通对象存储功能并且申请相关密钥即可使用上传下载
//----------------------------------------------------------

// COSPath 腾讯云COS地址
var COSPath = ""

// COSSecretId 腾讯云COSSecretId
var COSSecretId = ""

// COSSecretKey 腾讯云COSSecretKey
var COSSecretKey = ""

//-------------------------
// ApplePassword 苹果共享密钥
//-------------------------

var ApplePassword = ""

//-------------------
// Azure 的chatGPT服务
//-------------------

var AzureKey = ""
var AzureURL = ""
