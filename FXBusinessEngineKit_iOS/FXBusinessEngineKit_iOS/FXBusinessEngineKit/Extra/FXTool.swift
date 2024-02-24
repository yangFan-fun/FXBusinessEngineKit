//
//  FXTool.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import UIKit
import CocoaLumberjack
import Photos


public func FXColor(_ hexColor: String) -> UIColor {
    var red: UInt64 = 0
    var green: UInt64 = 0
    var blue: UInt64 = 0
    // 移除开头的 #
    let hex = String(hexColor[hexColor.index(hexColor.startIndex, offsetBy: 1)...])
    
    Scanner(string: String(hex[..<hex.index(hex.startIndex, offsetBy: 2)])).scanHexInt64(&red)
    Scanner(string: String(hex[hex.index(hex.startIndex, offsetBy: 2)..<hex.index(hex.startIndex, offsetBy: 4)])).scanHexInt64(&green)
    Scanner(string: String(hex[hex.index(hex.startIndex, offsetBy: 4)...])).scanHexInt64(&blue)
    
    return UIColor(red: CGFloat(red) / 255.0, green: CGFloat(green) / 255.0, blue: CGFloat(blue) / 255.0, alpha: 1.0)
}


public func FXColor(_ hexColor: String, alpha: CGFloat) -> UIColor {
    var red: UInt64 = 0
    var green: UInt64 = 0
    var blue: UInt64 = 0
    // 移除开头的 #
    let hex = String(hexColor[hexColor.index(hexColor.startIndex, offsetBy: 1)...])
    
    Scanner(string: String(hex[..<hex.index(hex.startIndex, offsetBy: 2)])).scanHexInt64(&red)
    Scanner(string: String(hex[hex.index(hex.startIndex, offsetBy: 2)..<hex.index(hex.startIndex, offsetBy: 4)])).scanHexInt64(&green)
    Scanner(string: String(hex[hex.index(hex.startIndex, offsetBy: 4)...])).scanHexInt64(&blue)
    
    return UIColor(red: CGFloat(red) / 255.0, green: CGFloat(green) / 255.0, blue: CGFloat(blue) / 255.0, alpha: alpha)
}


// MARK: -


public func FXRegularFont(_ size: CGFloat) -> UIFont {
    let font = UIFont(name: "PingFangSC-Regular", size: size)
    guard let font = font else {
        return UIFont.systemFont(ofSize: size)
    }
    return font
}

func FXMediumFont(_ size: CGFloat) -> UIFont {
    let font = UIFont(name: "PingFangSC-Medium", size: size)
    guard let font = font else {
        return UIFont.systemFont(ofSize: size)
    }
    return font
}

func FXSemiboldFont(_ size: CGFloat) -> UIFont {
    let font = UIFont(name: "PingFangSC-Semibold", size: size)
    guard let font = font else {
        return UIFont.systemFont(ofSize: size)
    }
    return font
}

func FXBoldFont(_ size: CGFloat) -> UIFont {
    return UIFont.boldSystemFont(ofSize: size)
}


// MARK: -


func FXStatusHeight() -> CGFloat {
    let height = UIApplication.shared.statusBarFrame.size.height
    
    return height
}

func FXSafeBottomHeight() -> CGFloat {
    let height = FXCurrentViewController()?.view.safeAreaInsets.bottom
    return height ?? 0
}


// MARK: -


func FXLog(info: String) {
    let isShowable = FXBusinessEngineConfig.isLogShowable
    if isShowable == false {
        return
    }
    DDLogInfo(info)
}


func FXRequestPhotoLibraryAuthorization(complete: @escaping (Bool) -> ()) {
    if PHPhotoLibrary.authorizationStatus() == .authorized {
        complete(true)
    } else {
        PHPhotoLibrary.requestAuthorization { status in
            DispatchQueue.main.async {
                complete(status == .authorized)
            }
        }
    }
}
