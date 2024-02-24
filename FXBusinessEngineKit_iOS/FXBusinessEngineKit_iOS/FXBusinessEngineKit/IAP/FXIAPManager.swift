//
//  FXIAPManager.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/2/21.
//

import Foundation
import StoreKit


enum FXIAPStatus {
    case purchasing
    case purchased
    case failed
    case restored
    case deferred
}

class FXIAPModel: Codable {
    init() {}
    
    var userId: Int64?
    var uuid: String?
    var isVip: Int?
    var expires: Int64?
    var transactionId: String?
    var originalTransactionId: String?
    var goodsId: String?
    var productId: Int64?
}

class FXIAPManager: NSObject, SKPaymentTransactionObserver, SKProductsRequestDelegate {
    
    // 缓存的商品信息
    var productCached: [SKProduct]?
    
    var product: (([SKProduct]?, Error?) -> (Void))?
    var buyStatus: ((FXIAPStatus) -> (Void))?
    
    static let staticInstance = FXIAPManager()
    
    deinit {
        SKPaymentQueue.default().remove(self)
        NotificationCenter.default.removeObserver(self)
    }
    
    override init() {
        super.init()
        SKPaymentQueue.default().add(self)

    }
    
    private func isAllowBuy() -> Bool {
        if SKPaymentQueue.canMakePayments() == true {
            return true
        }
        return false
    }
    
    
    /// 获取商品信息
    /// - Parameters:
    ///   - productsId: 商品Id
    ///   - complete: 完成后的回调
    func fetchProduct(productsId: [String], complete: @escaping ([SKProduct]?, Error?) -> (Void)) {
        if isAllowBuy() == false {
            complete(nil, nil)
            return
        }
        
        self.product = complete
        
        var _productsId: Set<String> = []
        for item in productsId {
            _productsId.insert(item)
        }
        let request = SKProductsRequest(productIdentifiers: _productsId)
        request.delegate = self
        request.start()
        
        FXLog(info: "[IAP] 正在获取商品信息")
    }
    
    
    /// 开始交易
    /// - Parameters:
    ///   - product: 一个具体商品
    ///   - complete: 当前交易状态
    func startBuy(product: SKProduct, complete: @escaping (FXIAPStatus) -> (Void)) {
        
        // 开始交易
        let payment = SKPayment(product: product)
        SKPaymentQueue.default().add(payment)
        FXLog(info: "[IAP] 开始交易，正在拉起交易弹窗")
        self.buyStatus = complete
    }
    
    // 监听购买状态
    func paymentQueue(_ queue: SKPaymentQueue, updatedTransactions transactions: [SKPaymentTransaction]) {
        for transaction in transactions {
            switch transaction.transactionState {
            case .purchasing:
                self.buyStatus?(.purchasing)
                break
            case .purchased:
                serverCheck(transaction: transaction)
                break
            case .restored:
                finishedTransaction(transaction: transaction)
                self.buyStatus?(.restored)
                break
            case .failed:
                finishedTransaction(transaction: transaction)
                self.buyStatus?(.failed)
                break
            case .deferred:
                finishedTransaction(transaction: transaction)
                self.buyStatus?(.deferred)
                break
            @unknown default:
                break
            }
        }
    }
    
    private func serverCheck(transaction: SKPaymentTransaction) {
        let receipt = Bundle.main.appStoreReceiptURL
        guard let receipt = receipt else { return }
        let recepData = try? Data(contentsOf: receipt)
        guard let recepData = recepData else { return }
        let recepString = recepData.base64EncodedString()
        let dict = [
            "receiptData" : recepString,
            "transactionId" : transaction.transactionIdentifier ?? "",
            "productId" : transaction.payment.productIdentifier
        ]
        
        FXNetworkManager.request(url: "/goods/v1/verifyGood", method: .post, parameter: dict, decoder: FXIAPModel.self) { model, error in
            self.finishedTransaction(transaction: transaction)
            if model?.isVip == 1 {
                self.buyStatus?(.purchased)
                self.finishedTransaction(transaction: transaction)
                FXUserManager.getUserInfo()
            } else {
                // 订单校验失败
            }
        }
    }
    
    
    // MARK: -

    
    // 获取商品具体信息
    func productsRequest(_ request: SKProductsRequest, didReceive response: SKProductsResponse) {
        let product = response.products
        let current = product
        self.product?(current, nil)
        FXLog(info: "[IAP] 正在获取商品信息")
    }
    
    func requestDidFinish(_ request: SKRequest) {
        FXLog(info: "[IAP] 获取商品信息结束")
    }
    
    func request(_ request: SKRequest, didFailWithError error: Error) {
        FXLog(info: "[IAP] 获取商品信息失败")
        self.product?(nil, error)
    }
    
    // 订单结束
    private func finishedTransaction(transaction: SKPaymentTransaction) {
        SKPaymentQueue.default().finishTransaction(transaction)
    }
    
    
}
