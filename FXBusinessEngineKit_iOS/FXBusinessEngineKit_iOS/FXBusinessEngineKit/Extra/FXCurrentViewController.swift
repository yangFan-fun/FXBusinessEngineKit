//
//  FXCurrentViewController.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import UIKit


func backToRootViewController() {
    let vc = FXCurrentViewController()
    vc?.navigationController?.popToRootViewController(animated: true)
}

/// 获取最顶层的视图控制器
/// - Returns: 视图控制器
func FXCurrentViewController() -> (UIViewController?) {
   var window = UIApplication.shared.keyWindow
   if window?.windowLevel != UIWindow.Level.normal{
     let windows = UIApplication.shared.windows
     for  windowTemp in windows{
       if windowTemp.windowLevel == UIWindow.Level.normal{
          window = windowTemp
          break
        }
      }
    }
   let vc = window?.rootViewController
   return FXCurrentViewController(vc)
}


private func FXCurrentViewController(_ vc :UIViewController?) -> UIViewController? {
   if vc == nil {
      return nil
   }
   if let presentVC = vc?.presentedViewController {
      return FXCurrentViewController(presentVC)
   }
   else if let tabVC = vc as? UITabBarController {
      if let selectVC = tabVC.selectedViewController {
          return FXCurrentViewController(selectVC)
       }
       return nil
    }
    else if let naiVC = vc as? UINavigationController {
       return FXCurrentViewController(naiVC.visibleViewController)
    }
    else {
       return vc
    }
 }
