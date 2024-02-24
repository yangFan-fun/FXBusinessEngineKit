package goods

import (
	"FXBusinessEngineKit_Server/configuration"
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/tool"
	"FXBusinessEngineKit_Server/user/userController"
	"FXBusinessEngineKit_Server/user/userModel"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// status
// 21000 向App Store的请求未使用HTTP POST请求方法。
// 21001 App Store 不再发送此状态代码。
// 21002 属性中的数据receipt-data格式错误或服务遇到临时问题。再试一次。
// 21003 系统无法验证收据。
// 21004 您提供的共享密钥与您账户中存档的共享密钥不匹配。
// 21005 收据服务器暂时无法提供收据。再试一次。
// 21006 此收据有效，但订阅处于过期状态。当您的服务器收到此状态代码时，系统还会解码并返回收据数据作为响应的一部分。此状态仅针对自动续订订阅的 iOS 6 样式交易收据返回。
// 21007 该收据来自测试环境，但您将其发送到生产环境进行验证。
// 21008 该收据来自生产环境，但您将其发送到测试环境进行验证。
// 21009 内部数据访问错误。稍后再试。
// 21010 系统找不到用用户账户或用户账户已被删除。

type VerifyUrl string

const (
	// SandboxUrl 沙箱验证地址
	SandboxUrl VerifyUrl = "https://sandbox.itunes.apple.com/verifyReceipt"
	// ProductUrl 现网验证地址
	ProductUrl VerifyUrl = "https://buy.itunes.apple.com/verifyReceipt"
)

// AppleVerityResponseInApp 包含所有应用内购买交易的应用内购买收据字段的数组
type AppleVerityResponseInApp struct {
	CancellationDate        string `json:"cancellation_date"`          // App Store 退还交易或从家庭共享中撤销交易的时间，采用类似于 ISO 8601 的日期时间格式。此字段仅适用于已退款或撤销的交易。
	CancellationDateMs      string `json:"cancellation_date_ms"`       // App Store 退还交易或从家庭共享中撤销交易的时间，采用 UNIX 纪元时间格式，以毫秒为单位。此字段仅适用于退款或撤销的交易。使用此时间格式处理
	CancellationDatePst     string `json:"cancellation_date_pst"`      // App Store 退还交易或从家庭共享中撤销交易的时间，在太平洋时区。此字段仅适用于退款或撤销的交易。
	CancellationReason      string `json:"cancellation_reason"`        // 退款或撤销交易的原因。值“1”表示客户由于您的应用程序中的实际或感知问题而取消了他们的交易。值“0”表示交易因其他原因被取消；例如，如果客户意外购买。可能的值：1, 0
	ExpiresDate             string `json:"expires_date"`               // 订阅到期时间或续订时间，采用类似于 ISO 8601 的日期时间格式。
	ExpiresDateMs           string `json:"expires_date_ms"`            // 订阅到期或续订的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式处理日期
	ExpiresDatePst          string `json:"expires_date_pst"`           // 订阅到期或续订的时间（太平洋时区）。
	IsInIntroOfferPeriod    string `json:"is_in_intro_offer_period"`   // 指示自动续订订阅是否处于推介价格期的指标。请参阅获取更多信息。
	IsTrialPeriod           string `json:"is_trial_period"`            // 指示订阅是否处于免费试用期内。请参阅获取更多信息。
	OriginalPurchaseDate    string `json:"original_purchase_date"`     // 原始应用内购买的时间，采用类似于 ISO 8601 的日期时间格式。
	OriginalPurchaseDateMS  string `json:"original_purchase_date_ms"`  // 原始应用内购买的时间，采用 UNIX 纪元时间格式，以毫秒为单位。对于自动续订订阅，该值表示首次购买订阅的日期。原始购买日期适用于所有产品类型，并且在同一产品 ID 的所有交易中保持不变。该值对应于 StoreKit 中原始交易的属性。使用此时间格式来处理日期。
	OriginalPurchaseDatePst string `json:"original_purchase_date_pst"` // 原始应用内购买的时间（太平洋时区）。
	OriginalTransactionId   string `json:"original_transaction_id"`    // 原始购买的交易标识符。
	ProductId               string `json:"product_id"`                 // 所购买产品的唯一标识符。您在 App Store Connect 中创建产品时提供此值，它对应于存储在交易的 payment 属性中的对象的属性。
	PromotionalOfferId      string `json:"promotional_offer_id"`       // 用户兑换的订阅优惠的标识符
	PurchaseDate            string `json:"purchase_date"`              // App Store 向用户账户收取购买或恢复产品费用的时间，或者 App Store 向用户账户收取订阅购买或过期后续订费用的时间，采用类似于 ISO 8601 的日期时间格式。
	PurchaseDateMS          string `json:"purchase_date_ms"`           // 对于消耗性、非消耗性和非续订订阅产品，App Store 向用户账户收取购买或恢复的产品费用的时间，采用 UNIX 纪元时间格式，以毫秒为单位。对于自动续订订阅，App Store 在订阅购买或续订失效后向用户账户收取费用的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式来处理日期。
	PurchaseDatePst         string `json:"purchase_date_pst"`          // App Store 向用户账户收取购买或恢复产品费用的时间，或者 App Store 向用户账户收取订阅购买或过期后续订费用的时间（采用太平洋时区）。
	Quantity                string `json:"quantity"`                   // 购买的消耗品数量。此值对应于SKPayment存储在交易的支付属性中的对象的数量属性。“1”除非使用可变付款进行修改，否则该值通常是不变的。最大值为 10。
	TransactionId           string `json:"transaction_id"`             // 交易的唯一标识符，例如购买、恢复或续订
	WebOrderLineItemId      string `json:"web_order_line_item_id"`     // 跨设备购买事件的唯一标识符，包括订阅续订事件。该值是识别订阅购买的主键。
}

type VerifyAppleResponseLatest struct {
	CancellationDateMS    string `json:"cancellation_date_ms"` // 退款时间
	CancellationReason    string `json:"cancellation_reason"`
	OriginalTransactionId string `json:"original_transaction_id"`
	ProductId             string `json:"product_id"`
	ExpiresDateMS         string `json:"expires_date_ms"`
	TransactionId         string `json:"transaction_id"`
}

// VerifyApplePending 自动续订订阅的信息
type VerifyApplePending struct {
	AutoRenewProductId        string `json:"auto_renew_product_id"`         // 自动续订订阅的当前续订首选项。此键的值对应于客户订阅续订的产品的属性。
	AutoRenewStatus           string `json:"auto_renew_status"`             // 自动续订订阅的当前续订状态。有关详细信息，请参阅。auto_renew_status 可能的值：1, 0
	ExpirationIntent          string `json:"expiration_intent"`             // 订阅过期的原因。此字段仅适用于包含过期的自动续订订阅的收据。有关详细信息，请参阅。expiration_intent 可能的值：1, 2, 3, 4, 5
	GracePeriodExpiresDate    string `json:"grace_period_expires_date"`     // 订阅续订宽限期到期的时间，采用类似于 ISO 8601 的日期时间格式。
	GracePeriodExpiresDateMS  string `json:"grace_period_expires_date_ms"`  // 订阅续订宽限期到期的时间，采用 UNIX 纪元时间格式，以毫秒为单位。仅当应用程序启用了计费宽限期并且用户在续订时遇到计费错误时，此键才会出现。使用此时间格式来处理日期。
	GracePeriodExpiresDatePST string `json:"grace_period_expires_date_pst"` // 订阅续订宽限期到期的时间（太平洋时区）。
	IsInBillingRetryPeriod    string `json:"is_in_billing_retry_period"`    // 指示 Apple 正在尝试自动续订过期订阅的标志。仅当自动续订订阅处于计费重试状态时，此字段才会出现。有关详细信息，请参阅。is_in_billing_retry_period 可能的值：1, 0
	OfferCodeRefName          string `json:"offer_code_ref_name"`           // 您在 App Store Connect 中配置的订阅优惠的参考名称。当客户兑换订阅优惠代码时，会出现此字段。有关详细信息，请参阅。
	OriginalTransactionId     string `json:"original_transaction_id"`       // 原始购买的交易标识符。
	PriceConsentStatus        string `json:"price_consent_status"`          // 订阅价格上涨的价格同意状态。仅当 App Store 通知客户价格上涨时，此字段才会出现。如果客户同意，默认值 为"0"并更改为。"1" 可能的值：1, 0
	ProductId                 string `json:"product_id"`                    // 所购买产品的唯一标识符。您在 App Store Connect 中创建产品时提供此值，它对应于存储在交易属性中的对象的属性。productIdentifierSKPaymentpayment
	PromotionalOfferId        string `json:"promotional_offer_id"`          // 用户兑换的自动续订订阅的促销优惠的标识符。当您在 App Store Connect 中创建促销优惠时，您可以在促销优惠标识符字段中提供此值。有关详细信息，请参阅。promotional_offer_id
	PriceIncreaseStatus       string `json:"price_increase_status"`         // 可能的值：1, 0
}

// AppleVerifyResponse 苹果订单校验
// InApp 数组不按时间顺序排序，解析数组时，需要迭代所有元素以确保满足需求，比如说，不能判断数组中最后一个元素是最新的订单
// 对于自动续订的收据，检查 latest_receipt_info 字段查看续订状态
// 查看数组为空，代表没有任何订单
type AppleVerifyResponse struct {
	Status             int64                       `json:"status"`               // 如果0 收据有效，则返回状态代码；如果出现错误，则返回状态代码。状态代码反映了应用程序收据的
	Environment        string                      `json:"environment"`          // 生成收据的环境：Sandbox(沙盒), Production(正式环境)
	IsRetryable        string                      `json:"is_retryable"`         // 请求期间发生错误时的指示器。值1表示暂时问题；稍后重试验证此收据。值0表示存在无法解决的问题；不要重新尝试验证此收据。这仅适用于状态代码21100–21199。
	LatestReceipt      string                      `json:"latest_receipt"`       // 最新的 Base64 编码的应用收据。这仅返回包含自动续订订阅的收据。
	LatestReceiptInfo  []VerifyAppleResponseLatest `json:"latest_receipt_info"`  // 包含所有应用内购买交易的数组。这不包括您的应用标记为已完成的消费品的交易。
	PendingRenewalInfo []VerifyApplePending        `json:"pending_renewal_info"` // 在 JSON 文件中，一个数组，其中每个元素包含所标识的每个自动续订订阅的待续订信息。这只返回包含自动续订订阅的应用程序收据。product_id
	Receipt            struct {
		AdamId                     int64  `json:"adam_id"`                      //见。app_item_id
		AppItemId                  int64  `json:"app_item_id"`                  // 由 App Store Connect 生成并由 App Store 用于唯一标识所购买的应用程序。仅在生产中为应用分配此标识符。将此值视为 64 位长整数。
		ApplicationVersion         string `json:"application_version"`          // 应用程序的版本号。应用程序的版本号对应于. 在生产中，此值是设备上基于. 在沙盒中，该值始终为。CFBundleVersionCFBundleShortVersionStringInfo.plist receipt_creation_date_ms"1.0"
		BundleId                   string `json:"bundle_id"`                    // 收据所属应用的捆绑包标识符
		DownloadId                 int64  `json:"download_id"`                  // 应用下载交易的唯一标识符。
		ExpirationDate             string `json:"expiration_date"`              // 通过批量购买计划购买的应用程序的收据到期时间，采用类似于 ISO 8601 的日期时间格式。
		ExpirationDateMs           string `json:"expiration_date_ms"`           // 通过批量购买计划购买的应用程序的收据到期时间，采用 UNIX 纪元时间格式，以毫秒为单位。如果通过批量购买计划购买的应用程序没有此密钥，则收据不会过期。使用此时间格式处理日期。
		ExpirationDatePst          string `json:"expiration_date_pst"`          // 通过批量购买计划购买的应用程序的收据过期时间（太平洋时区）
		OriginalApplicationVersion string `json:"original_application_version"` // 用户最初购买的应用程序的版本。该值不会改变，对应于原始购买文件中的（在 iOS 中）或String（在 macOS 中）的值。在沙盒环境中，该值始终为。CFBundleVersionCFBundleShortVersionInfo.plist"1.0"
		OriginalPurchaseDate       string `json:"original_purchase_date"`       // 原始应用购买的时间，采用类似于 ISO 8601 的日期时间格式。
		OriginalPurchaseDateMs     string `json:"original_purchase_date_ms"`    // 原始应用购买的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式处理日期。
		OriginalPurchaseDatePst    string `json:"original_purchase_date_pst"`   // 原始应用程序购买时间（太平洋时区）
		PreorderDate               string `json:"preorder_date"`                // 用户订购可预购应用的时间，采用类似于 ISO 8601 的日期时间格式。
		PreorderDataMS             string `json:"preorder_data_ms"`             // 用户订购可供预订的应用程序的时间，采用 UNIX 纪元时间格式，以毫秒为单位。仅当用户预购应用程序时，此字段才会出现。使用此时间格式来处理日期。
		PreorderDatePst            string `json:"preorder_date_pst"`            // 用户订购可供预订的应用程序的时间（太平洋时区）。
		ReceiptCreationDate        string `json:"receipt_creation_date"`        // App Store 生成收据的时间，采用类似于 ISO 8601 的日期时间格式。
		ReceiptCreationDateMs      string `json:"receipt_creation_date_ms"`     // App Store 生成收据的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式处理日期。这个值不会改变。
		ReceiptCreationDatePst     string `json:"receipt_creation_date_pst"`    // App Store 生成收据的时间（太平洋时区）。
		ReceiptType                string `json:"receipt_type"`                 // 生成的收据类型。该值对应于应用程序或 VPP 购买的环境。 可能的值：Production, ProductionVPP, ProductionSandbox, ProductionVPPSandbox
		RequestDate                string `json:"request_date"`                 // 处理对端点的请求并生成响应的时间，采用类似于 ISO 8601 的日期时间格式。verifyReceipt
		RequestDateMs              string `json:"request_date_ms"`              // 处理对端点的请求并生成响应的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式处理日期。verifyReceipt
		RequestDatePst             string `json:"request_date_pst"`             // 在太平洋时区处理对端点的请求并生成响应的时间。verifyReceipt
		VersionExternalIdentifier  int64  `json:"version_external_identifier"`  // 标识应用程序修订版的任意数字。在沙箱中，该键的值为“0”。

		InApp []AppleVerityResponseInApp `json:"in_app"` //包含所有应用内购买交易的应用内购买收据字段的数组
	} `json:"receipt"` //发送以供验证的收据的 JSON 表示形式。
}

// AppleVerifyTransactionNotification 苹果服务端订单通知
type AppleVerifyTransactionNotification struct {
	AutoRenewAdamId              string                                `json:"auto_renew_adam_id"`                // App Store Connect 生成且 App Store 用于唯一标识用户订阅续订的自动续订订阅的标识符。将此值视为 64 位整数。
	AutoRenewProductId           string                                `json:"auto_renew_product_id"`             // 用户订阅续订的自动续订订阅的产品标识符。
	AutoRenewStatus              string                                `json:"auto_renew_status"`                 // 自动续订订阅产品的当前续订状态。请注意，这些值是字符串“true”和，而不是收据中出现的“false”字符串“0”或。“1”auto_renew_status
	AutoRenewStatusChangeDate    string                                `json:"auto_renew_status_change_date"`     // 用户打开或关闭自动续订订阅的续订状态的时间，采用类似于 ISO 8601 标准的日期时间格式。
	AutoRenewStatusChangeDateMS  string                                `json:"auto_renew_status_change_date_ms"`  // 用户打开或关闭自动续订订阅的续订状态的时间（UNIX 时间），以毫秒为单位。使用此时间格式来处理日期。
	AutoRenewStatusChangeDatePST string                                `json:"auto_renew_status_change_date_pst"` // 用户打开或关闭自动续订订阅的续订状态的时间（太平洋标准时间）。
	Bid                          string                                `json:"bid"`                               // 包含应用程序包 ID 的字符串。
	Bvrs                         string                                `json:"bvrs"`                              // 包含应用程序包版本的字符串。
	Deprecation                  string                                `json:"deprecation"`                       // App Store 服务器通知 V1被弃用的日期。有关更多信息，请参阅App Store 服务器通知变更日志。
	Environment                  string                                `json:"environment"`                       // App Store 生成收据的环境。 可能的值：Sandbox, PROD
	NotificationType             string                                `json:"notification_type"`                 // 触发通知的订阅事件。请参阅 参考资料 获取版本 1 通知类型的完整列表。
	OriginalTransactionId        int64                                 `json:"original_transaction_id"`           // App Store 服务器通知发送通知的原始事务标识符。
	Password                     string                                `json:"password"`                          // 与验证收据时在passwordApp Store 收据字段中提交的共享密钥的值相同。
	UnifiedReceipt               VerifyAppleNotificationUnifiedReceipt `json:"unified_receipt"`                   // 包含有关应用程序最近的应用内购买交易的信息的对象。
}

// VerifyAppleNotificationUnifiedReceipt 苹果服务端通知统一的应用内购买对象
type VerifyAppleNotificationUnifiedReceipt struct {
	Environment        string                      `json:"environment"`          // App Store 生成收据的环境。 可能的值：Sandbox, Production
	LatestReceipt      string                      `json:"latest_receipt"`       // 最新的 Base64 编码的应用收据。
	LatestReceiptInfo  []VerifyAppleResponseLatest `json:"latest_receipt_info"`  // 包含解码值的最新 100 个应用内购买交易的数组。此数组不包括您的应用程序标记为已完成的消费品的交易。该数组的内容与用于收据验证的verifyReceipt端点响应中的内容相同。latest_receiptresponseBody.Latest_receipt_info
	PendingRenewalInfo []VerifyApplePending        `json:"pending_renewal_info"` // 一个数组，其中每个元素都包含 中标识的每个自动续订订阅的待续订信息。该数组的内容与用于收据验证的verifyReceipt端点响应中的内容相同。product_idresponseBody.Pending_renewal_info
	Status             int                         `json:"status"`               // 状态代码，0表示通知有效。
}

type VerifyParameter struct {
	Receipt  string `json:"receipt-data"`
	Password string `json:"password"`
	IsNew    bool   `json:"exclude-old-transactions"`
}

type VerifyClientParameter struct {
	Receipt       string `json:"receiptData"`   // 凭证
	ProductId     string `json:"productId"`     // 商品id
	TransactionId string `json:"transactionId"` // 交易id
}

// RequestUrlTimeout 请求超时时间
var requestUrlTimeout = 15 * time.Second

// password 苹果共享密钥
var password = configuration.ApplePassword

// AppleVerifyReceipt 苹果订单验证
func AppleVerifyReceipt(url VerifyUrl, transactionModel VerifyClientParameter, context *gin.Context) (AppleVerifyResponse, error) {
	userId, _ := context.Get("userId")
	userIdStr := fmt.Sprintf("%s", userId)

	var passwordString string
	var productId = tool.ProductIdFormHeader(context)

	// 可以验证多个产品及id
	if productId == 124567 {
		passwordString = password
	}

	parameter := VerifyParameter{
		Receipt:  transactionModel.Receipt,
		Password: passwordString,
		IsNew:    true,
	}

	// 将参数序列化成json
	parameterJson, jsonErr := json.Marshal(parameter)
	if jsonErr != nil {
		log.RecordLog(log.Err, "[苹果订单校验]数据转换失败")
		return AppleVerifyResponse{}, jsonErr
	}
	// 发起请求开始校验
	client := http.Client{
		Timeout: requestUrlTimeout,
	}
	response, err := client.Post(string(url), "application/json", bytes.NewBuffer(parameterJson))
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[苹果订单校验]订单校验请求发起失败%s, userId %s", err, userIdStr))
		return AppleVerifyResponse{}, err
	}
	if response.StatusCode != 200 {
		log.RecordLog(log.Err, "[苹果订单校验]苹果订单校验发起请求连接失败")
		responseNetLinkErr := errors.New("网络连接失败")
		return AppleVerifyResponse{}, responseNetLinkErr
	}

	// 将苹果返回信息转成字符串
	strBuffer := new(bytes.Buffer)
	_, _ = strBuffer.ReadFrom(response.Body)
	bodyStr := strBuffer.String()
	//fmt.Printf("订单信息 %s", bodyStr)

	defer response.Body.Close()

	// 将数据解析成模型
	var responseModel AppleVerifyResponse
	re := json.Unmarshal([]byte(bodyStr), &responseModel)
	if re != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[苹果订单校验]数据转成模型失败%s, userId %s", re, userIdStr))
		return AppleVerifyResponse{}, re
	}

	// 订单校验状态码异常
	if responseModel.Status != 0 {
		log.RecordLog(log.Err, fmt.Sprintf("[苹果订单校验]订单校验异常，苹果返回错误码 %d, userId %s", responseModel.Status, userIdStr))
		transactionErr := errors.New(fmt.Sprintf("订单校验异常，苹果返回错误码 %d", responseModel.Status))
		return responseModel, transactionErr
	}

	// 开始校验订单信息
	verifyErr := VerifyResponseSuccessUserUpdatable(responseModel, context)
	if verifyErr != nil {
		return responseModel, verifyErr
	}

	// 所有校验已完成，没有异常
	return responseModel, nil
}

// VerifyResponseSuccessUserUpdatable 订单验证成功的话存储用户信息
func VerifyResponseSuccessUserUpdatable(response AppleVerifyResponse, context *gin.Context) error {
	// 将参数绑定到模型
	var transactionParameter VerifyClientParameter
	parameterError := context.ShouldBindJSON(&transactionParameter)
	if parameterError != nil {
		return parameterError
	}

	userId := tool.UserIdFromHeader(context)
	productId := tool.ProductIdFormHeader(context)

	for _, value := range response.Receipt.InApp {
		// 根据客户端当前交易id进行匹配订单
		if transactionParameter.TransactionId == value.TransactionId {

			var expireTime, _ = strconv.ParseInt(value.ExpiresDateMs, 10, 64)
			var isVip = 0
			if expireTime >= time.Now().UnixMilli() {
				isVip = 1
			}

			ne := expireTime / 1000
			ut := userModel.FXUserAdapterWalletModel{
				TransactionId:         value.TransactionId,
				OriginalTransactionId: value.OriginalTransactionId,
				ProductId:             value.ProductId,
				ExpirationDate:        ne,
				IsVip:                 isVip,
			}

			// 入库
			err := userController.UserSaveOrder(userId, productId, ut)
			if err != nil {
				return err
			}
			return nil
		}
	}

	log.RecordLog(log.Err, fmt.Sprintf("[苹果订单验证]用户信息入库失败，当前交易没有匹配成功, userId %d", userId))
	saveErr := errors.New("用户信息入库失败，当前交易没有匹配成功")
	return saveErr
}

func AppleVerifyNotification(context *gin.Context) {
	// 将请求体中的数据绑定到模型
	var model AppleVerifyTransactionNotification
	err := context.ShouldBindJSON(&model)
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[苹果交易苹果服务端通知]数据绑定失败 %s", err))
		return
	}
	// 获取当前苹果推送的订单id
	currentTransactionId := model.OriginalTransactionId
	currentTransactionIdStr := strconv.FormatInt(currentTransactionId, 10)

	// 在统一收据对象中根据当前订单id获取具体的订单信息
	var isVip int
	var expiresDate int64
	var goodsId string
	var transactionId string
	for _, receipt := range model.UnifiedReceipt.LatestReceiptInfo {
		if receipt.OriginalTransactionId == currentTransactionIdStr {
			// 获取过期时间
			expiresDateM, _ := strconv.ParseInt(receipt.ExpiresDateMS, 10, 64)
			expiresDate = expiresDateM / 1000
			goodsId = receipt.ProductId
			transactionId = receipt.TransactionId
		}
	}
	for _, order := range model.UnifiedReceipt.PendingRenewalInfo {
		if order.OriginalTransactionId == currentTransactionIdStr {
			// 是否正在订阅中
			isVipInt, _ := strconv.Atoi(order.AutoRenewStatus)
			isVip = isVipInt
		}
	}

	userId := tool.UserIdFromHeader(context)
	productId := tool.ProductIdFormHeader(context)

	date := time.Now().Format(configuration.DateFormat)

	w := userModel.FXUserWallet{
		UserId:                userId,
		IsVip:                 isVip,
		Expires:               expiresDate,
		TransactionId:         transactionId,
		OriginalTransactionId: currentTransactionIdStr,
		GoodsId:               goodsId,
		ProductId:             productId,
		AppleNotification:     1,
		TransactionCreateTime: date,
	}

	// 从数据库中找到这个用户，设置订阅信息
	userModel.GormUpdateUserWithTransaction(w)

}
