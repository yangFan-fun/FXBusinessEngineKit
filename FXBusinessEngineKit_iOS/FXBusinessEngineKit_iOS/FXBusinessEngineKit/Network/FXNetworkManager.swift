//
//  FXNetwork.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import UIKit
import Alamofire
import CoreTelephony
import CoreServices
import SnapKit


enum FaceArtNetworkStatus {
    case unReachable
    case viaWiFi
    case viaWWAN
    case unInitialize
}


class FXNetworkManager {
    
    struct FXNetworkUploadModel: Codable {
        var url: String?
    }
    
    static let shared = FXNetworkManager()
    
    // 网络管理器
    private var networkReachability: NetworkReachabilityManager? = NetworkReachabilityManager()
    // 网络状态
    var networkStatus: FaceArtNetworkStatus = .unInitialize
    // 网络状态空页面
    var networkEmptyView: UIView?

    // 域名
    let domainURL = FXBusinessEngineConfig.domainPath
    
    let kNetworkStatusChange = "kNetworkStatusChange"
    
    init() {
        _initNetwork()
        _initializeNetwokReachability()
    }
    
    /// 发起网络请求
    /// - Parameters:
    ///   - url:        路由地址
    ///   - method:     请求方式
    ///   - parameter:  参数
    ///   - decoder:    要序列化的模型对象类型
    ///   - complete:   完成后的回调
    static func request<T: Codable>(serverType: FXBusinessServerType = .base, url: String, method: HTTPMethod, parameter: [String : Any]?, decoder: T.Type, complete: @escaping (T?, Error?) -> (Void)) {
        let u = self.shared._requestPath(type: serverType) + url
        AF.request(u, method: method, parameters: parameter, encoding: JSONEncoding.default, headers: self.shared._headers(), interceptor: nil, requestModifier: nil).response { response in
            switch response.result {
            case .success(let data):
                self.shared.serializationData(data: data, decoder: decoder, complete: complete)
            case .failure(let error):
                complete(nil, error)
            }
        }
    }
    
    
    /// 上传图片
    /// - Parameters:
    ///   - data:       图片数据
    ///   - progress:   上传进度
    ///   - complete:   完成后的回调
    static func uploadImage(data: Data, progress: ((Progress) -> ())? = nil, complete: @escaping (FXNetworkUploadModel?, Error?) -> (Void)) {
        let u = self.shared.domainURL + "/upload/v1/image"
        AF.upload(multipartFormData: { part in
            let imageName = FXBusinessEngineConfig.imageGeneralNameRule
            part.append(data, withName: "file", fileName: imageName, mimeType: "image/jpg")
        }, to: u, method: .post, headers: self.shared._headers()).uploadProgress { _progress in
            progress?(_progress)
        }.responseData { response in
            switch response.result {
            case .success(let data):
                self.shared.serializationData(data: data, decoder: FXNetworkUploadModel.self, complete: complete)
            case .failure(let error):
                complete(nil, error)
            }
        }
    }
    
    
    /// 上传视频
    /// - Parameters:
    ///   - path:       视频地址
    ///   - progress:   上传进度
    ///   - complete:   完成后的回调
    static func uploadVideo(path: URL, progress: ((Progress) -> ())? = nil, complete: @escaping (FXNetworkUploadModel?, Error?) -> (Void)) {
        let u = self.shared.domainURL + "/upload/v1/video"
        AF.upload(multipartFormData: { part in
            let videoName = FXBusinessEngineConfig.videoGeneralNameRule
            part.append(path, withName: "file", fileName: videoName, mimeType: "video/mp4")
        }, to: u, method: .post, headers: self.shared._headers()).uploadProgress { _progress in
            progress?(_progress)
        }.responseData { response in
            switch response.result {
            case .success(let data):
                self.shared.serializationData(data: data, decoder: FXNetworkUploadModel.self, complete: complete)
                break
            case .failure(let error):
                complete(nil, error)
                break
            }
        }
    }
    
    
    /// 序列化网络请求后的数据
    /// - Parameters:
    ///   - data:       源数据
    ///   - decoder:    要序列化的模型对象类型
    ///   - complete:   完成后的回调
    func serializationData<T: Codable>(data: Data?, decoder: T.Type, complete: @escaping (T?, Error?) -> ()) {
        guard let data = data else {
            let error = NSError(domain: "接口返回数据为空", code: -1, userInfo: nil)
            complete(nil, error)
            return
        }
        let dict = try? JSONSerialization.jsonObject(with: data) as? NSDictionary
        guard let dict = dict else {
            let error = NSError(domain: "序列化成JSON失败", code: -1, userInfo: nil)
            complete(nil, error)
            return
        }
        let dictData = dict["data"] as? String
        guard let dictData = dictData else {
            let error = NSError(domain: "获取Data字段值失败", code: -1, userInfo: nil)
            complete(nil, error)
            return
        }
        let str = dictData.data(using: .utf8, allowLossyConversion: false)
        guard let str = str else {
            let e = NSError(domain: "Data字段转换成Data类型失败", code: -1, userInfo: nil)
            complete(nil, e)
            return
        }
        let model = try? JSONDecoder().decode(decoder.self, from: str)
        complete(model, nil)
    }
    
    static func _destination(_ targetPath: URL,
                             _ response: HTTPURLResponse) -> (URL, DownloadRequest.Options) {
        return (targetPath, [.removePreviousFile, .createIntermediateDirectories])
    }
    
    static let destination = DownloadRequest.suggestedDownloadDestination(for: .documentDirectory)
    
    static func download(_ url: String, complete: @escaping (URL?, Error?, Float?) -> (Void)) {
        let videoPath = FXSandBox.videoCached.path + "/" + FXBusinessEngineConfig.videoGeneralNameRule
        AF.download(url, to: { return self._destination(URL(string: videoPath)!, $1)}
        ).downloadProgress { progress in
            complete(nil, nil, Float(progress.fractionCompleted))
        }.validate().response { response in
            switch response.result {
            case .success(let _url):
                complete(_url, nil, nil)
            case .failure(let error):
                complete(nil, error, nil)
            }
        }
    }
    
    static func downloadImage(_ url: String, name: String, complete: @escaping (String?, Error?, Float?) -> (Void)) {
        let imagePath = FXSandBox.imageCached.path + "/" + FXBusinessEngineConfig.imageGeneralNameRule
        
        let path = imagePath
        AF.download(url, to: { return self._destination(URL(fileURLWithPath: path), $1)}
        ).downloadProgress { progress in
            complete(nil, nil, Float(progress.fractionCompleted))
        }.validate().response { response in
            switch response.result {
            case .success(_):
                complete(imagePath, nil, nil)
            case .failure(let error):
                complete(nil, error, nil)
            }
        }
    }
    
    func _initNetwork() {
        let conf = URLSessionConfiguration.af.default
        conf.timeoutIntervalForRequest = FXBusinessEngineConfig.timeoutForRequest
        let _ = Alamofire.Session(configuration: conf)
    }
    
    
    // MARK: -
    
    
    private func _headers() -> HTTPHeaders {
        let uuid = FXKeyChain.shared.getUUID()
        let userId = FXUserManager.loginModel?.userId
        let token = FXUserManager.loginModel?.token
        var dict = [
            "Content-Type" : "application/json",
            "platform" : "iOS",
            "uuid" : uuid,
            "productId" : FXBusinessEngineConfig.productId,
        ] as HTTPHeaders
        
        if let token = token {
            dict["token"] = "Bearer\(token)"
        }
        if let userId = userId {
            dict["userId"] = "\(userId)"
        }
        return dict
    }
    
    
    // MARK: -
    
    
    // 网络状态管理
    private func _initializeNetwokReachability() {
        networkReachability?.startListening(onQueue: .main, onUpdatePerforming: { result in
            
            switch result {
            case.notReachable, .unknown:
                self._showNetworkEmptyViewIfNeeded()
                self.networkStatus = .unReachable
                FXLog(info: "[\(#function)] 网络类型：无网络或者无网络权限")
            case.reachable(.ethernetOrWiFi):
                self._hideNetworkEmptyViewIfNeeded()
                self.networkStatus = .viaWiFi
                FXLog(info: "[\(#function)] 网络类型：WiFi")
            case .reachable(.cellular):
                self._hideNetworkEmptyViewIfNeeded()
                self.networkStatus = .viaWWAN
                FXLog(info: "[\(#function)] 网络类型：蜂窝网络")
            }
            NotificationCenter.default.post(name: Notification.Name(self.kNetworkStatusChange), object: nil)
        })
    }
    
    private func _showNetworkEmptyViewIfNeeded() {
        let view = FXNetworkEmptyView()
        self.networkEmptyView = view
        let currentViewController = FXCurrentViewController()
        guard let currentViewController = currentViewController else {
            return
        }

        currentViewController.view.addSubview(view)
        view.snp.makeConstraints { make in
            make.edges.equalTo(currentViewController.view)
        }
    }
    
    private func _hideNetworkEmptyViewIfNeeded() {
        self.networkEmptyView?.removeFromSuperview()
    }
    
    func _requestPath(type: FXBusinessServerType) -> String {
        var url: String = ""
        switch type {
        case .base:
            url = FXBusinessEngineConfig.domainPath + ":" + String(FXBusinessEngineConfig.basePort)
            break
        case .inference:
            url = FXBusinessEngineConfig.inferenceDomainPath + ":" + String(FXBusinessEngineConfig.inferencePort)
            break
        }
        return url
    }
}
