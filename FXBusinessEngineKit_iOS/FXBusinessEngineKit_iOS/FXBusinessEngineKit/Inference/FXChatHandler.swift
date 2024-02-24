//
//  FXChatHandler.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/2/24.
//

import Foundation

struct FXChatModel: Codable {
    /// 消息内容
    var string: String
}

class FXChatHandler {
    
    static let shared = FXChatHandler()
    
    private init() {}
    
    
    /// 发送一条消息
    /// - Parameters:
    ///   - content: 消息内容
    ///   - complete: 完成后的回调
    static func sendMessage(content: String, complete: @escaping (FXChatModel?, Error?) -> ()) {
        let parameter = [
            "content" : content
        ]
        FXNetworkManager.request(serverType: .inference, url: "/gpt/v1/chat", method: .post, parameter: parameter, decoder: FXChatModel.self, complete: complete)
    }
    
}
