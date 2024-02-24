//
//  FXToast.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import Foundation
import Toast_Swift


class FXToast {
    
    static let staticInstance = FXToast()
    
    static func sharedInstance() -> FXToast {
        return staticInstance
    }
    
    init() {
        
    }
    
    static func showToast(_ text: String?, _ view: UIView?) {
        guard let view = view else { return }
        view.makeToast(text, duration: 4, position: .center)
    }
    
}
