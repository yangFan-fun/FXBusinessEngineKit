//
//  FXChatGPT.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/2/22.
//

import Foundation


struct FXChatGPTMessageModel: Codable {
    
    /// 消息内容
    var content: String
}

class FXChatGPTHandler {
    
    let shared = FXChatGPTHandler()
    
    private init() {}
    
    
    /// 发送一条对话消息
    /// - Parameters:
    ///   - content: 消息内容
    ///   - complete: 完成后的回调
    static func sendMessage(content: String, complete: @escaping (FXChatGPTMessageModel?, Error?) -> ()) {
        let parameter = [
            "content" : content
        ]
        FXNetworkManager.request(url: "/gpt/v1/chat", method: .post, parameter: parameter, decoder: FXChatGPTMessageModel.self, complete: complete)
    }
}
