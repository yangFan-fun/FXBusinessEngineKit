//
//  FXRootViewController.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import UIKit
import MobileCoreServices


class FXRootViewController: UIViewController, UINavigationControllerDelegate, UIImagePickerControllerDelegate {
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        _setupUI()
    }
    
    func _setupUI() {
        view.backgroundColor = UIColor.white
        
        FXUserManager.login { model in
            let a = model
        }
        
//        DispatchQueue.main.asyncAfter(deadline: .now() + 5) {
//            let image = UIImage(named: "test")?.jpegData(compressionQuality: 1)
//            guard let image = image else { return }
//            FXNetworkManager.upload(url: "/uploadOutsea/v1/face", image: image) { model, error in
//                if error != nil {
//                    return
//                }
//                
//            }
//        }
        
        let picker = UIImagePickerController()
        picker.sourceType = .photoLibrary
        picker.mediaTypes = [kUTTypeMovie as String]
        picker.delegate = self
        self.present(picker, animated: true)
    }
    
    func imagePickerControllerDidCancel(_ picker: UIImagePickerController) {
        
    }
    
    func imagePickerController(_ picker: UIImagePickerController, didFinishPickingMediaWithInfo info: [UIImagePickerController.InfoKey : Any]) {
        
        let mediaType = info[UIImagePickerController.InfoKey.mediaType] as? String
        if mediaType == kUTTypeMovie as String {
            let url = info[UIImagePickerController.InfoKey.mediaURL] as? URL
            print("视频地址：\(url)")
            
            guard let url = url else { return }
            
            FXNetworkManager.uploadVideo(path: url) { progress in
                let _progress = progress.completedUnitCount
                let _total = progress.totalUnitCount
                let _curr = progress.completedUnitCount
                let _p = _curr / _total
                print("上传的进度：\(_p) 进度数据：\(progress.fractionCompleted)")
            } complete: { model, error in
                print("上传后的地址：\(model?.url)")
            }

        }
        
    }
}
