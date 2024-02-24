//
//  AppDelegate.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import UIKit


@main
class AppDelegate: UIResponder, UIApplicationDelegate {

    var window: UIWindow?
    
    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        
        window = UIWindow(frame: UIScreen.main.bounds)
        window?.rootViewController = UINavigationController(rootViewController: FXRootViewController())
        window?.makeKeyAndVisible()
        
        _initBusinessEngineConfig()
        
        return true
    }
    
    func _initBusinessEngineConfig() {
        FXBusinessEngineConfig.setDomain(path: "https://www.lipsticked.cloud:8080")
//        FXBusinessEngineConfig.setImageNameRuleForUpload(rule: "1.jpeg")
    }
}

