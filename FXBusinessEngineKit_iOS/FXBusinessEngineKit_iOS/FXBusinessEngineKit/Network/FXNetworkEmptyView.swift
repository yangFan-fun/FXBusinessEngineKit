//
//  FXNetworkEmptyView.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import UIKit


class FXNetworkEmptyView: UIView {
    override init(frame: CGRect) {
        super.init(frame: frame)
        
        _initializeUI()
    }
    
    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
    
    func _initializeUI() {
        self.backgroundColor = FXColor("#FFFFFF")
        
        let imageView = UIImageView()
        imageView.image = UIImage(named: "network")
        self.addSubview(imageView)
        imageView.snp.makeConstraints { make in
            
            make.width.equalTo(200)
            make.height.equalTo(200)
            make.centerX.equalTo(self)
            make.centerY.equalTo(self).offset(-100)
        }
        
        let title = UILabel()
        title.text = "网络无法连接，请打开网络权限或检查网络".local()
        title.textColor = FXColor("#000000")
        title.font = FXMediumFont(17)
        title.textAlignment = .center
        self.addSubview(title)
        title.snp.makeConstraints { make in
            make.top.equalTo(imageView.snp.bottom).offset(100)
            make.left.equalTo(20)
            make.right.equalTo(-20)
        }
    }
}
