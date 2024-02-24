//
//  FXUser.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import Foundation
import KeychainAccess


/// 用户登录模型
struct FXUserLoginModel: Codable {
    var token: String?
    var uuid: String?
    var userId: Int?
}

/// 用户信息模型
struct FXUserModel: Codable {
    var userId: Int?
    var uuid: String?
    var registrationDate: Double?
    var wallet: FXUserWalletModel?
    var productId: Int?
}

/// 用户钱包模型
struct FXUserWalletModel: Codable {
    var userId: Int?
    var uuid: String?
    var isVip: Int? // 0 没有开通会员 1 已经开通会员
    var expires: Double?
    var transactionId: String?
    var originalTransactionId: String?
    var goodsId: String?
    var productId: Int?
}


class FXUserManager {
    
    static let shared = FXUserManager()
    
    private init() {
        
    }
    
    private let kUserCached = "FXUserModelCached"
    
    /// 用户信息
    static var user: FXUserModel? {
        get {
            return _userInfoGetable()
        }
    }
    
    /// 用户id
    static var userId: Int {
        get {
            let _userId = self.user?.userId
            guard let _userId = _userId else {
                FXLog(info: "\(#function) 用户信息获取失败，请重新登录")
                return 0
            }
            return _userId
        }
    }
    
    /// 是否是会员
    static var isVip: Bool {
        let _isVip = user?.wallet?.isVip
        if _isVip == nil || _isVip == 0 {
            return false
        }
        return true
    }
    
    /// 会员过期时间
    static var vipExpires: String? {
        let stamp = user?.wallet?.expires
        guard let stamp = stamp else { return nil }
        let date = FXLocalizationDateWithTemplate(stamp: stamp, template: "yyyy-MM-dd")
        return date
    }
    
    /// 用户登录数据
    /// 将 {token}以及 {userId} 放在网络请求的header中，才可以成功发起网络请求
    static var loginModel: FXUserLoginModel?
    
    
    /// 用户登录，首次登录注册
    /// 成功登录后会自动获取一次用户信息
    /// 需要在应用启动时最先调用当前方法以获取token用来发起网络请求
    static func login() {
        FXNetworkManager.request(url: "/user/v1/login", method: .post, parameter: nil, decoder: FXUserLoginModel.self) { model, error in
            if error != nil {
                FXLog(info: "\(#function) fail")
                return
            }
            self.loginModel = model
            getUserInfo()
        }
    }
    
    
    /// 用户登录，首次登录注册
    /// 需要在应用启动时最先调用当前方法以获取token用来发起网络请求
    /// - Parameter complete: 是否登录成功状态
    static func login(complete: @escaping (Bool) -> ()) {
        FXNetworkManager.request(url: "/user/v1/login", method: .post, parameter: nil, decoder: FXUserLoginModel.self) { model, error in
            if error != nil {
                FXLog(info: "\(#function) fail")
                complete(false)
                return
            }
            self.loginModel = model
            complete(true)
        }
    }
    
    
    /// 获取用户信息 仅缓存用户数据
    static func getUserInfo() {
        FXNetworkManager.request(url: "/user/v1/user", method: .post, parameter: nil, decoder: FXUserModel.self) { model, error in
            if error != nil {
                FXLog(info: "\(#function) fail")
                return
            }
            if model == nil {
                FXLog(info: "\(#function) fail")
                return
            }
            self._userModelCachabel(model: model)
        }
    }
    
    
    /// 获取用户信息
    /// - Parameter complete: 获取用户信息
    static func getUserInfo(complete: @escaping (FXUserModel?, Error?) -> ()) {
        FXNetworkManager.request(url: "/user/v1/user", method: .post, parameter: nil, decoder: FXUserModel.self) { model, error in
            if error != nil {
                FXLog(info: "\(#function) fail")
                complete(nil, error)
                return
            }
            self._userModelCachabel(model: model)
            complete(model, nil)
        }
    }
    
    
    /// 用户信息缓存
    /// 在获取用户信息时无需实时调用接口
    /// - Parameter model: 缓存的用户信息
    static private func _userModelCachabel(model: FXUserModel?) {
        guard let model = model else { return }
        let data = try? JSONEncoder().encode(model)
        guard let data = data else { return }
        UserDefaults.standard.set(data, forKey: FXUserManager.shared.kUserCached)
        UserDefaults.standard.synchronize()
    }
    
    
    /// 获取用户信息缓存
    /// - Returns: 被缓存的用户信息
    static private func _userInfoGetable() -> FXUserModel? {
        let data = UserDefaults.standard.object(forKey: FXUserManager.shared.kUserCached) as? Data
        guard let data = data else { return nil }
        let model = try? JSONDecoder().decode(FXUserModel.self, from: data)
        return model
    }
    
}
