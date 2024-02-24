package main

import (
	"FXBusinessEngineKit_Server/configuration"
	"FXBusinessEngineKit_Server/goods"
	"FXBusinessEngineKit_Server/gpt/gptController"
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/middleware"
	"FXBusinessEngineKit_Server/tool"
	"FXBusinessEngineKit_Server/upload/uploadController"
	"FXBusinessEngineKit_Server/user/userController"
	"fmt"
	"github.com/gin-gonic/gin"
)

//----------------------------------------------------
// FXBusinessEngineKit_Server
// 拆分了两部分服务分别是
// 当前的 FXBusinessEngineKit_Server（基础服务）
// 还有文件夹中已经编译好的 NaturalBridge Server (推理服务)
//----------------------------------------------------

//----------------------------------------------------
// FXBusinessEngineKit_Server 功能列表
// 1. 请求校验
// 2. 用户体系
// 3. 对象存储
// 4. GPT 由 NaturalBridge 服务实现
// 5. iOS订单校验
// 6. SSL

// NaturalBridge 功能列表
// 1. 图片换脸
// 2. 视频换脸
//----------------------------------------------------

//----------------------------------------------------
// 应用部署环境：
// CentOS 7
// Nginx > 1.18
// MySQL > 5.7.28
// Gorm
//----------------------------------------------------

//----------------------------------------------------
// ⚠️⚠️⚠️ 使用之前需要先查看 configuration
//----------------------------------------------------

func main() {
	log.InitLog()
	tool.InitConnectDatabase()

	startServer()
}

func startServer() {
	ginServer := gin.Default()
	ginServer.Use(middleware.LoadTLS())

	connectAPI := ginServer.Group("/connect")
	{
		connectAPI.GET("/base", func(context *gin.Context) {
			tool.ResponseSuccess(context, "FXBusinessEngineKit_Server请求测试成功")
		})
	}

	// 用户体系相关路由
	userGroup := ginServer.Group("/user")
	{
		// 用户登录
		// 应用启动时第一个调用的接口，用户获取鉴权Token
		userGroup.POST("/v1/login", middleware.JWTUserLogin(), middleware.LogMiddleware(), userController.UserLogin)

		// 获取用户信息
		userGroup.GET("/v1/user", middleware.JWTMiddleware(), middleware.LogMiddleware(), userController.User)
	}

	// 商品相关路由
	goodsGroup := ginServer.GET("/goods")
	{
		// 获取商品列表
		goodsGroup.GET("/v1/fetchGoods", middleware.JWTMiddleware(), middleware.LogMiddleware(), goods.FetchGoods)

		// 苹果订单验证
		goodsGroup.POST("/v1/iOS/verifyReceipt", middleware.JWTMiddleware(), middleware.LogMiddleware(), goods.VerifyReceipt)
	}

	// 上传相关路由
	uploadGroup := ginServer.Group("/upload")
	{
		// 上传图片
		uploadGroup.POST("/v1/image", middleware.JWTMiddleware(), middleware.LogMiddleware(), uploadController.UploadFace)
		// 上传视频
		uploadGroup.POST("/v1/video", middleware.JWTMiddleware(), middleware.LogMiddleware(), uploadController.UploadVideo)
	}

	// GPT相关路由
	// 由NaturalBridge 服务实现
	gptGroup := ginServer.Group("/gpt")
	{
		// 聊天，服务供应商：Azure
		gptGroup.POST("/v1/chat", middleware.JWTMiddleware(), middleware.LogMiddleware(), gptController.Chat)
	}

	// 开启Https
	err := ginServer.RunTLS(configuration.AppPort, configuration.HttpsCertFile, configuration.HttpsKeyFile)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("Https服务开启失败：%s", err))
		return
	}

	// 不开启Https
	//err := ginServer.Run(configuration.AppPort)
	//if err != nil {
	//	log.RecordLog(log.Err, fmt.Sprintf("Http服务开启失败：%s", err))
	//	return
	//}
}
