//
//  FXKeyChain.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import Foundation
import KeychainAccess


class FXKeyChain {
    
    static let shared = FXKeyChain()
    
    let kuuid = "uuid"
    
    let kService = "keychainService"
    
    let keychain: Keychain!
    
    private init() {
        keychain = Keychain(service: kService)
        _cachedUUID()
    }
    
    func _cachedUUID() {
        
        let uuid = UUID().uuidString
        
        let cachedUUID = keychain[kuuid]
        if cachedUUID != nil {
            return
        }
        
        keychain[kuuid] = uuid
    }
    
    func getUUID() -> String {
        let uuid = keychain[kuuid]
        return uuid ?? ""
    }
}
