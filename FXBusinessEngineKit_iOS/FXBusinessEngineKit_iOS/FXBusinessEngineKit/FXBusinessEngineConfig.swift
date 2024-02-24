//
//  FXBusinessEngineConfig.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import Foundation


enum FXBusinessServerType {
    case base
    case inference
}

class FXBusinessEngineConfig {
    
    static let shared = FXBusinessEngineConfig()
    
    /// 获取域名
    static var domainPath: String { get { self.shared.domainPath } }
    
    /// 获取基础服务端口号
    /// 如果不设置默认为：8081
    static var basePort: Int { get { self.shared.basePort } }
    
    /// 获取当前应用的productId
    static var productId: String { get { self.shared.productId } }
    
    /// 是否打印日志，默认打印日志
    static var isLogShowable: Bool { get { self.shared.logShowable } }
    
    /// 网络请求超时时间
    /// 默认时间 30s
    static var timeoutForRequest: TimeInterval { get { self.shared.timeoutForRequest } }
    
    /// 图片的命名规则
    static var imageGeneralNameRule: String { get { self.shared.fetchImageNameRule() } }
    
    /// 视频的命名规则
    static var videoGeneralNameRule: String { get { self.shared.fetchVideoNameRule() } }
    
    /// 是否开启推理服务
    static var isOpenInferenceServer: Bool { get { self.shared.isOpenInferenceServer } }
    
    /// 获取推理服务域名
    static var inferenceDomainPath: String { get { self.shared.fetchInferenceDomainPath() } }
    
    /// 获取推理服务端口号
    /// 如果不设置默认为：8082
    static var inferencePort: Int { get { self.shared.inferencePort } }
    
    private var domainPath: String = ""
    private var basePort: Int = 8081
    private var productId: String = ""
    private var logShowable: Bool = true
    private var timeoutForRequest: TimeInterval = 30
    private var imageNameRule: String = ""
    private var videoNameRule: String = ""
    private var isOpenInferenceServer: Bool = false
    private var inferenceDomainPath: String = ""
    private var inferencePort: Int = 8082
    
    private init() {}
    
    
    // MARK: -
    
    
    /// 设置域名
    static func setDomain(path: String) {
        self.shared.domainPath = path
    }
    
    /// 设置基础服务端口号
    /// 不设置默认为：8081
    static func setBasePort(port: Int) {
        self.shared.basePort = port
    }
    
    /// 设置当前应用的productId，用于服务端调用接口鉴权以及用量统计
    static func setProductId(_ id: String) {
        self.shared.productId = id
    }
    
    /// 设置是否打印日志
    static func setLogShowable(allow: Bool) {
        self.shared.logShowable = allow
    }
    
    /// 设置网络请求超时时间
    static func setNetworkTimeoutIntervalForRequest(second: TimeInterval) {
        self.shared.timeoutForRequest = second
    }
    
    /// 设置图片命名规则
    /// 下载，上传默认使用此规则
    static func setImageGeneralNameRule(rule: String) {
        self.shared.imageNameRule = rule
    }
    
    /// 设置视频命名规则
    /// 下载，上传默认使用此规则
    static func setVideoGeneralNameRule(rule: String) {
        self.shared.videoNameRule = rule
    }
    
    /// 设置是否开启推理服务
    static func isOpenInferenceServer(open: Bool) {
        self.shared.isOpenInferenceServer = open
    }
    
    /// 设置推理服务域名
    static func setInferenceServerDomainPath(path: String) {
        self.shared.inferenceDomainPath = path
    }
    
    /// 设置推理服务端口号
    /// 不设置默认为：8082
    static func setInferenceServerPort(port: Int) {
        self.shared.inferencePort = port
    }
    
    
    // MARK: -
    
    
    // 默认的上传下载图片命名规则
    private func normalImageNameRule() -> String {
        let second = currentTimeStamp()
        let userId = FXUserManager.userId
        let imageName = "image_\(second)_\(userId).jpeg"
        return imageName
    }
    
    // 上传下载图片命名规则
    private func fetchImageNameRule() -> String {
        if self.imageNameRule.isEmpty == false {
            return self.imageNameRule
        }
        return normalImageNameRule()
    }
    
    // 默认的上传下载视频命名规则
    private func normalVideoNameRule() -> String {
        let second = currentTimeStamp()
        let userId = FXUserManager.userId
        let videoName = "video_\(second)_\(userId).mp4"
        return videoName
    }
    
    // 上传下载视频命名规则
    private func fetchVideoNameRule() -> String {
        if self.videoNameRule.isEmpty == false {
            return self.videoNameRule
        }
        return normalVideoNameRule()
    }
    
    // 获取推理服务的域名
    func fetchInferenceDomainPath() -> String {
        if self.inferenceDomainPath.isEmpty == true {
            return self.domainPath
        }
        return self.inferenceDomainPath
    }
    
    
    // MARK: -
    
    
    private func currentTimeStamp() -> Int {
        let date = Date().timeIntervalSince1970
        let second = Int(date)
        return second
    }
}
