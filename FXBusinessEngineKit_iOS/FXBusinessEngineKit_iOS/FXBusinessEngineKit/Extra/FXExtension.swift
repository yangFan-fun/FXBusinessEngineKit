//
//  FXExtension.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import UIKit
import SDWebImage


extension UIImageView {
    
    func FXSetImage(url: String?) {
        guard let url = url else { return }
        let urlString = URL(string: url)
        guard let urlString = urlString else { return }
        self.sd_setImage(with: urlString, placeholderImage: UIImage(named: "imagePlaceholder"))
    }
}


extension UIImage {
    func  FXResizeImage(reSize: CGSize)-> UIImage  {
        UIGraphicsBeginImageContextWithOptions (reSize, false , UIScreen.main.scale);
        self.draw(in: CGRect(x: 0, y: 0, width: reSize.width, height: reSize.height))
        let reSizeImage: UIImage = UIGraphicsGetImageFromCurrentImageContext () ?? UIImage();
        UIGraphicsEndImageContext ();
        return  reSizeImage;
    }
}


extension String {
    
    /// 本地化文案
    /// - Returns: 本地化文案
    func local() -> String {
        let str = NSLocalizedString(self, comment: "")
        return str
    }
    
    /// 获取文字高度
    /// - Parameters:
    ///   - width: 固定宽度值
    ///   - font: 字体
    /// - Returns: 高度
    func FXStringHeight(width: CGFloat, font: UIFont) -> CGFloat {
        let widthMax = ceil(width)
        let widthInt = Int(widthMax)
        let size = self.boundingRect(with: CGSize(width: widthInt, height: .max), options: .usesLineFragmentOrigin, attributes: [.font : font], context: nil).size
        let height = ceil(size.height)
        return height
    }
    
    
    /// 获取文字高度，宽度为屏幕宽度
    /// - Parameter font: 字体
    /// - Returns: 文字高度
    func FXStringHeight(font: UIFont) -> CGFloat {
        let height = FXStringHeight(width: UIScreen.main.bounds.size.width, font: font)
        return height
    }
}


extension UIView {
    
    /// 添加渐变图层
    /// - Parameters:
    ///   - startPoint: 开始坐标
    ///   - endPoint: 结束坐标
    ///   - colors: 颜色
    ///   - locations: 位置
    @objc func gradientColor(_ startPoint: CGPoint, _ endPoint: CGPoint, _ colors: [Any],_ locations:[NSNumber]) {
        
        guard startPoint.x >= 0, startPoint.x <= 1, startPoint.y >= 0, startPoint.y <= 1, endPoint.x >= 0, endPoint.x <= 1, endPoint.y >= 0, endPoint.y <= 1 else {
            return
        }
        _ = delay(0, task: {[self] in
            // 强制刷新一次
            layoutIfNeeded()
            
            var gradientLayer: CAGradientLayer!
            removeGradientLayer()
            gradientLayer = CAGradientLayer()
            gradientLayer.locations = locations
            gradientLayer.frame = self.layer.bounds
            gradientLayer.startPoint = startPoint
            gradientLayer.endPoint = endPoint
            gradientLayer.colors = colors
            gradientLayer.cornerRadius = self.layer.cornerRadius
            gradientLayer.masksToBounds = true
            self.layer.insertSublayer(gradientLayer, at: 0)
            self.backgroundColor = UIColor.clear
            self.layer.masksToBounds = false
        })
    }
    
    @objc func removeGradientLayer() {
        if let sl = self.layer.sublayers {
            for layer in sl {
                if layer.isKind(of: CAGradientLayer.self) {
                    layer.removeFromSuperlayer()
                }
            }
        }
    }
}


extension UILabel{
    
    enum gradientDirection{
        case horizontal
        case vertical
    }
    
    /// 设置渐变图层
    /// - Parameters:
    ///   - dir: 方向
    ///   - colors: 颜色
    ///   - locations: 位置
    func gradientColor(dir:gradientDirection,_ colors: [Any],_ locations:[CGFloat]) {
        _ = delay(0, task: {
            let size = self.bounds.size
            UIGraphicsBeginImageContextWithOptions(size, false, UIScreen.main.scale)
            guard let context = UIGraphicsGetCurrentContext() else { return }
            let colorSpace = CGColorSpaceCreateDeviceRGB()
            ///设置渐变颜色
            let gradientRef = CGGradient(colorsSpace: colorSpace, colors: colors as CFArray, locations: locations)!
            var startPoint = CGPoint(x: size.width / 2, y: 0)
            var endPoint = CGPoint(x: size.width / 2, y: size.height)
            if dir == .horizontal {
                startPoint = CGPoint(x: 0, y: size.height/2)
                endPoint = CGPoint(x: size.width, y: size.height/2)
            }
            context.drawLinearGradient(gradientRef, start: startPoint, end: endPoint, options: CGGradientDrawingOptions(arrayLiteral: .drawsBeforeStartLocation,.drawsAfterEndLocation))
            let gradientImage = UIGraphicsGetImageFromCurrentImageContext()
            UIGraphicsEndImageContext()
            self.textColor = UIColor(patternImage: gradientImage!)
        })
    }
}


extension UIButton{
    enum gradientDir{
        case horizontal
        case vertical
    }
    func gradientColor(dir:gradientDir,_ colors: [Any],_ locations:[CGFloat],_ state:UIControl.State) {
        
        if self.bounds.size.width == 0 && self.bounds.size.height == 0 {
            return
        }
        
        _ = delay(0, task: {
            let size = self.bounds.size
            UIGraphicsBeginImageContextWithOptions(size, false, UIScreen.main.scale)
            guard let context = UIGraphicsGetCurrentContext() else{return}
            let colorSpace = CGColorSpaceCreateDeviceRGB()
            ///设置渐变颜色
            let gradientRef = CGGradient(colorsSpace: colorSpace, colors: colors as CFArray, locations: locations)!
            var startPoint = CGPoint(x: size.width / 2, y: 0)
            var endPoint = CGPoint(x: size.width / 2, y: size.height)
            if dir == .horizontal {
                startPoint = CGPoint(x: 0, y: size.height/2)
                endPoint = CGPoint(x: size.width, y: size.height/2)
            }
            context.drawLinearGradient(gradientRef, start: startPoint, end: endPoint, options: CGGradientDrawingOptions(arrayLiteral: .drawsBeforeStartLocation,.drawsAfterEndLocation))
            let gradientImage = UIGraphicsGetImageFromCurrentImageContext()
            UIGraphicsEndImageContext()
            self.setBackgroundImage(gradientImage, for: state)
        })
    }
}


typealias delayTask = (_ cancel : Bool) -> Void

func delay(_ time: TimeInterval, task: @escaping ()->()) -> delayTask? {

    func dispatch_later(block: @escaping ()->()) {
        let t = DispatchTime.now() + time
        DispatchQueue.main.asyncAfter(deadline: t, execute: block)
    }
    
    var closure: (()->Void)? = task
    var result: delayTask?

    let delayedClosure: delayTask = {
        cancel in
      if let internalClosure = closure {
            if (cancel == false) {
                DispatchQueue.main.async(execute: internalClosure)
            }
        }
        closure = nil
        result = nil
    }

    result = delayedClosure

    dispatch_later {
        if let delayedClosure = result {
            delayedClosure(false)
        }
    }
    
  return result
}

func cancel(_ task: delayTask?) {
    task?(true)
}
