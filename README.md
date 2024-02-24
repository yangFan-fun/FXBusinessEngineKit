# FXBusinessEngineKit







#### 服务端部署环境：

FXBusinessEngineKit\_Server：

CentOS 7.0

NaturalBridge：

Ubuntu 20.04.6 LTS



***



#### 📌📌📌点击：[视频换脸，图片说话效果展示](https://gaomi-star.feishu.cn/docx/MPLid2qMAoyaI1xTVKLc5lxPnvf#part-BxIwddHTqoisJzxeJZVcLsRJnNd)



## 一、服务端API介绍

###### 使用前的准备工作

1. 数据库安装：MySQL > 5.7.28

2. 用户体系数据库建表

3. 应用详情数据库建表



Tips: 所有的服务端API无需关心返回值，客户端已经全部处理好，你需要做的是看清楚注释内容，是否满足你的需求，以及如何扩展功能。



## FXBusinessEngineKit\_Server API

> FXBusinessEngineKit\_Server 是服务端的基础服务，主要用来解决绑定识别用户设备，用户体系建立以及商业化相关等内容

#### 用户体系相关内容

```Go
// 用户体系相关路由
userGroup := ginServer.Group("/user")
{
    // 用户登录，没有登录过就注册当前用户
    // 应用启动时第一个调用的接口，用户获取鉴权Token，客户端请求API时服务端会校验此Token，通过后才会进行方法调用，Token有效期默认30天
    // 工具类产品一般没有手机验证码或者谷歌以及邮箱注册，简单的做法是客户端会拿到设备的UUID(设备唯一标识)将此UUID存储在应用内，每次请求API的时候传入此UUID
    // 就可以做到识别一个用户或者绑定一部设备
    userGroup.POST("/v1/login", middleware.JWTUserLogin(), middleware.LogMiddleware(), userController.UserLogin)

    // 获取用户信息
    // 通过登录接口获取Token之后，客户端请求当前接口来获取用户信息，在使用前的准备工作中有两种和用户体系相关的数据库表
    // 这个接口就是查询在数据库中用户信息，包含了：
    
    // 基础信息：
    // 1.UUID 用户的设备唯一标识
    // 2.用户的注册时间
    // 3.当前用户所在组织机构或者在那个产品下注册的
    // 4.平台信息
    
    // 钱包信息：
    // 1.是否是会员
    // 2.会员过期时间
    // 3.当前开通会员的订单id
    // 4.当前开通会员的苹果原始订单id
    // 5.钱包创建时间
    // 6.用户的金币余额
    userGroup.GET("/v1/user", middleware.JWTMiddleware(), middleware.LogMiddleware(), userController.User)
}
```



#### 订单校验 - iOS

```Go
// 商品相关路由
goodsGroup := ginServer.GET("/goods")
{
    // 获取商品列表
    // 如果是在应用运营期间要动态并且不发版的改变应用内的商品信息，可以通过当前API来设置返回给客户端的商品，让客户端在每次应用启动时获取一次商品
    goodsGroup.GET("/v1/fetchGoods", middleware.JWTMiddleware(), middleware.LogMiddleware(), goods.FetchGoods)

    // 苹果订单验证
    // 需要将苹果的共享密钥设置在配置文件中
    goodsGroup.POST("/v1/iOS/verifyReceipt", middleware.JWTMiddleware(), middleware.LogMiddleware(), goods.VerifyReceipt)
}
```

苹果IAP（应用内购买）简介：

1. **苹果应用商店商品类型**

2. **苹果应用内购买的收入**

3. **站在服务端的角度上看待苹果的IAP**

4. **苹果应用内购买整体流程**

5. **苹果服务端通知类型详细介绍**

6. **苹果订单校验苹果返回JSON示例**



#### 上传到COS对象存储 - 腾讯云

```Go
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

//----------------------------------------------------------

// 上传相关路由
uploadGroup := ginServer.Group("/upload")
{
    // 上传图片
    // 上传是先保存到服务器中，然后上传到腾讯云的对象存储中
    // 返回一个链接
    uploadGroup.POST("/v1/image", middleware.JWTMiddleware(), middleware.LogMiddleware(), uploadController.UploadFace)
    // 上传视频
    // 上传是先保存到服务器中，然后上传到腾讯云的对象存储中
    // 返回
    uploadGroup.POST("/v1/video", middleware.JWTMiddleware(), middleware.LogMiddleware(), uploadController.UploadVideo)
}
```



## NaturalBridge API

> NaturalBridge 是服务端的推理服务，主要是封装了一些开源库用来实现一些AIGC的能力，主要解决的是在不用购买第三方SDK的情况下就可以实现类似的效果快速实现功能上线验证产品

#### 效果预览：

1. 图片换脸

2. 视频换脸

3. 图片中的人物说话



> 如果需要换脸功能，需要先在服务端部署好相应的项目环境
>
> 图片换脸使用的算法项目是：Roop
>
> 视频换脸使用的算法项目是：SimSwap
>
> NaturalBridge服务已经编译好，可以直接在当前算法项目的根目录使用 `nohup ./NaturalBridge &` 运行



> NaturalBridge中的换脸工作原理：
>
> ***
>
> Tips: 部署项目之后，NaturalBridge会生成以下内容：
>
> FXServer\_NaturalBridge：资源文件夹，换脸的源图片和目标图片或者视频存储在当前文件夹中
>
> FXBusinessEngineKit\_NaturalBridge\_record.log：日志文件
>
> FXServer\_NaturalBridge\_Roop\_RunShell.sh：调用图片换脸算法的脚本文件
>
> FXServer\_NaturalBridge\_SimSwap\_RunShell.sh：调用视频视频换脸算法的脚本文件
>
> ***
>
> 接收到换脸请求后会先从COS对象存储下载所有的素材信息，存储在素材目录中
>
> NaturalBridge在首次启动时会创建调用换脸算法的脚本文件，下载好资源素材后通过文件路径调用，脚本文件
>
> 启动算法后：
>
> 1. 会将当前任务标记为 Processing 状态，给当前任务创建一个TaskId，将用户的UserId和当前任务绑定
>
> 2. 为了最大程度的快速完成任务，一次只能处理一个任务，不能并发处理，在当前状态下如果有其他发起任务的请求会直接拒绝，并且返回一个 Waitting 的状态
>
> 3. NaturalBridge会监控算法输出，算法成功输出文件后当前任务会标记为 Processed 状态，等待接收其他任务
>
> 4. 获取到算法输出的文件后上传到COS对象存储中，并且返回一个地址



#### 图片换脸 - Roop

```Go
faceSwapGroup.POST("/v1/roopPhotoFaceSwap", Controller2.RoopPhotoFaceSwap)
```



#### 视频换脸 - SimSwap

```Go
faceSwapGroup.POST("/v1/simSwapVideoFaceSwap", Controller2.SimSwapVideoFaceSwap)
```



#### 照片说话 - SadTalker

```Go
faceSwapGroup.POST("/v1/sadTalkerPhotoTalk", Controller2.sadTalkerPhotoTalk)
```

#### 聊天 - OpenAI

```Go
gptGroup.POST("/v1/chat", chatController.Chat)
```



####

***



## 二、移动端 iOS API介绍：

#### 用户体系

```Swift
/// 用户登录，首次登录注册
/// 成功登录后会自动获取一次用户信息
/// 需要在应用启动时最先调用当前方法以获取token用来发起网络请求
static func login() {}


/// 获取用户信息
/// 开通会员或者需要在刷新用户信息时调用
static func getUserInfo() {}
static func getUserInfo(complete: @escaping (FXUserModel?, Error?) -> ()) {}
```



#### 订单校验 - iOS

```Swift
/// 获取商品信息
/// - Parameters:
///   - productsId: 商品Id
///   - complete: 完成后的回调
/// 在 App Store 配置商品后会有一个商品id，通过商品id调用当前方法可以获取一个 SKProduct 类型对象
/// 这个对象表示了一个商品的具体信息，比如说国际化后的价格信息，国际化后的商品描述
/// 在用户点击界面上的购买，在逻辑上开始交易后，拉起苹果的支付弹窗时，就是通过当前的 SKProduct 对象进行下单
func fetchProduct(productsId: [String], complete: @escaping ([SKProduct]?, Error?) -> (Void)) {}


/// 开始交易
/// - Parameters:
///   - product: 一个商品对象
///   - complete: 当前交易状态
/// 方法会回调交易的状态，分别有：

/// purchasing
/// purchased
/// failed
/// restored
/// deferred

/// 有两个维度的交易完成，第一个是苹果完成扣费，方法会回调 purchased 状态，但当前并没有给用户发货
/// 在 purchased 状态后会方法内部会调用 func serverCheck(transaction: SKPaymentTransaction)
/// 进行服务端订单校验，客户端会发送 receiptData，transactionId，productId 和 UserId
/// 前两个字段用于服务端校验当前订单是否有效，有效就发货，productId 和 UserId 用于将当前订单和当前用户绑定
func startBuy(product: SKProduct, complete: @escaping (FXIAPStatus) -> (Void)) {}
```



#### 上传到COS对象存储 - 腾讯云

```Swift
struct FXNetworkUploadModel: Codable {
/// 上传到COS对象存储后的访问链接
    var url: String?
}

/// 上传图片
/// - Parameters:
///   - data:       图片数据
///   - progress:   上传进度
///   - complete:   完成后的回调
static func uploadImage(data: Data, progress: ((Progress) -> ())? = nil, complete: @escaping (FXNetworkUploadModel?, Error?) -> (Void)) {}


/// 上传视频
/// - Parameters:
///   - path:       视频地址
///   - progress:   上传进度
///   - complete:   完成后的回调
static func uploadVideo(path: URL, progress: ((Progress) -> ())? = nil, complete: @escaping (FXNetworkUploadModel?, Error?) -> (Void)) {}
```



#### 图片换脸

```Swift
```



#### 视频换脸

```Swift
/// 视频换脸的任务有两个步骤
/// 1.提交一个视频换脸任务
/// 2.等待默认的代理通知，或者可以自己在业务层进行查询任务状态


/// 提交一个视频换脸任务
/// - Parameters:
///   - originSource: 一张要换脸的图片素材地址
///   - targetSource: 换脸的视频素材地址
///   - isNeedDelegateNotification: 是否需要开启代理通知，默认开启，任务出现异常或者任务执行结束后会通过代理方法
///    func faceSwap(handler: FXFaceSwapHandler, currentTask detail: FXFaceSwapTaskModel?, error: Error?) 来通知业务层
///    也可以不开启通知，业务层自行查询任务状态 查询路由：/faceSwap/v1/querySimSwapVideoFace 或者查询方法：queryTaskDetail
static func submitTask(originSource: String, targetSource: String, _ isNeedDelegateNotification: Bool = true) {}



/// 查询当前正在进行的换脸任务
/// - Parameter complete: 查询完成后的回调
static func queryTaskDetail(complete: @escaping (FXFaceSwapTaskModel?, Error?) -> ()) {}
```



#### 照片说话 - SadTalker

```Swift
```



#### 聊天 - OpenAI

```Swift
/// 发送一条消息
/// - Parameters:
///   - content: 消息内容
///   - complete: 完成后的回调
static func sendMessage(content: String, complete: @escaping (FXChatModel?, Error?) -> ()) {}
```

