//
//  FXFaceSwap.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/2/13.
//

import Foundation


enum FXFaceSwapTaskStatus: Codable {
    case processing
    case processed
    case waiting
}

@objc class FXFaceSwapTaskModel: NSObject, Codable {
    /// 当前的任务Id
    var taskId: String?
    /// 当前的任务状态
    var status: FXFaceSwapTaskStatus?
    /// 当前任务发起的用户Id
    var userId: Int?
    /// 当前的任务处理成功之后上传到腾讯云COS的地址
    var resultResource: String?
}


@objc protocol FXFaceSwapDelegate {
    
    /// 当前换脸任务的信息
    /// - Parameters:
    ///   - handler:    换脸处理器
    ///   - detail:     任务的具体信息
    ///   - error:      错误信息
    func faceSwap(handler: FXFaceSwapHandler, currentTask detail: FXFaceSwapTaskModel?, error: Error?)
}

class FXFaceSwapHandler: NSObject {
    
    static let shared = FXFaceSwapHandler()
    
    weak var delegate: FXFaceSwapDelegate?
    
    var queryTaskTimer: Timer?
    
    private override init() {}
    
    
    /// 提交一个视频换脸任务
    /// - Parameters:
    ///   - originSource: 一张要换脸的图片素材地址
    ///   - targetSource: 换脸的视频素材地址
    ///   - isNeedDelegateNotification: 是否需要开启代理通知，默认开启，任务出现异常或者任务执行结束后会通过代理方法
    ///    func faceSwap(handler: FXFaceSwapHandler, currentTask detail: FXFaceSwapTaskModel?, error: Error?) 来通知业务层
    ///    也可以不开启通知，业务层自行查询任务状态 查询路由：/faceSwap/v1/querySimSwapVideoFace 或者查询方法：queryTaskDetail
    static func submitTask(originSource: String, targetSource: String, _ isNeedDelegateNotification: Bool = true) {
        let parameter = [
            "originSource" : originSource,
            "targetSource" : targetSource
        ]
        FXNetworkManager.request(serverType: .inference, url: "/faceSwap/v1/simSwapVideoFaceSwap", method: .post, parameter: parameter, decoder: FXFaceSwapTaskModel.self) { model, error in
            
        }
        if isNeedDelegateNotification == false {
            return
        }
        DispatchQueue.main.asyncAfter(wallDeadline: .now() + 1) {
            self.shared.queryTaskDetail()
        }
    }
    
    
    /// 查询当前正在进行的换脸任务
    /// - Parameter complete: 查询完成后的回调
    static func queryTaskDetail(complete: @escaping (FXFaceSwapTaskModel?, Error?) -> ()) {
        FXNetworkManager.request(serverType: .inference, url: "/faceSwap/v1/querySimSwapVideoFace", method: .get, parameter: nil, decoder: FXFaceSwapTaskModel.self, complete: complete)
    }
    
    
    // MARK: -
    
    
    /// 查询换脸任务
    @objc private func queryTaskDetail() {
        queryTaskTimerInvalidate()
        FXNetworkManager.request(serverType: .inference, url: "/faceSwap/v1/querySimSwapVideoFace", method: .get, parameter: nil, decoder: FXFaceSwapTaskModel.self) { model, error in
            if error != nil {
                self.delegate?.faceSwap(handler: self, currentTask: model, error: error)
                return
            }
            
            // 当前任务是否完成需要判断
            // 1. 任务状态是否是完成状态
            // 2. 当前完成的任务是否是当前用户发起的
            if model?.status == .processed {
                let currentUser = FXUserManager.userId
                if model?.userId == currentUser {
                    self.delegate?.faceSwap(handler: self, currentTask: model, error: nil)
                    return
                }
            }
            
            let timer = Timer(timeInterval: 1, target: self, selector: #selector(self.queryTaskDetail), userInfo: nil, repeats: false)
            RunLoop.current.add(timer, forMode: .common)
            self.queryTaskTimer = timer
        }
    }
    
    private func queryTaskTimerInvalidate() {
        queryTaskTimer?.invalidate()
        queryTaskTimer = nil
    }
}
