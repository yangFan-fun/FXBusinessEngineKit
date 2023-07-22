Pod::Spec.new do |spec|

  spec.name         = "FXBusinessEngineKit"
  spec.version      = "1.0.0"
  spec.summary      = "FanXing-BusinessEngineKit"

  spec.description  = <<-DESC
  A commercial suite of tool based products that quickly realizes basic services, including user system, order verification service, dialogue service, cloud storage service, and inference service.
                   DESC

  spec.homepage     = "https://github.com/yangFan-fun/FXBusinessEngineKit"
  
  spec.license      = { :type => "MIT", :file => "LICENSE" }

  spec.author       = { "yangfan" => "DBasdyangfan@outlook.com" }
 
  spec.platform     = :ios

  spec.ios.deployment_target = "13.0"

  spec.source       = { :git => "https://github.com/yangFan-fun/FXBusinessEngineKit.git", :tag => "#{spec.version}" }

  # spec.source_files  = "Classes", "FXBusinessEngineKit/**/*.swift"

  spec.vendored_frameworks = "FXBusinessEngineKit.framework"

  spec.dependency "Alamofire", "~>5.0.0"

  spec.swift_version = "5.0"

end
