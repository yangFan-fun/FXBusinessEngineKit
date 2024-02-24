//
//  FXSandBox.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import Foundation


enum FXSandBox {
    
    /// Documents
    case home
    
    /// 日志文件
    case log
    
    /// 缓存的图片
    case imageCached
    
    /// 缓存的视频
    case videoCached
    
    /// 获取目录地址，没有则创建
    var path: String {
        var path: String
        switch self {
        case .home:
            path = NSHomeDirectory() + "/Documents"
        case .log:
            path = FXSandBox.home.path + "/Log"
        case .imageCached:
            path = FXSandBox.home.path + "/ImageCached"
        case .videoCached:
            path = FXSandBox.home.path + "/VideoPath"
        }
        
        var isDirectory: ObjCBool = false
        if !FileManager.default.fileExists(atPath: path, isDirectory: &isDirectory) || !isDirectory.boolValue {
            try? FileManager.default.removeItem(atPath: path)
            try? FileManager.default.createDirectory(atPath: path, withIntermediateDirectories: true)
        }
        return path
    }
        
}

/// 文件是否存在
/// - Parameter path: 路径
/// - Returns: 是否存在
func FXFileIsExists(path: String) -> Bool {
    let isExists = FileManager.default.fileExists(atPath: path)
    return isExists
}
