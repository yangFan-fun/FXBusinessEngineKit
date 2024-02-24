//
//  FXLocalization.swift
//  FXBusinessEngineKit
//
//  Created by 杨帆 on 2024/1/31.
//

import Foundation


enum FXLanguage {
    /// 汉语
    case zh
    /// 英语
    case en
    /// 法语
    case fr
    /// 西班牙语
    case es
    /// 意大利语
    case it
    /// 印尼语
    case id
    /// 德语
    case de
    /// 韩语
    case ko
    /// 土耳其语
    case tr
    /// 巴西葡萄牙语
    case pt_BR
    /// 未知
    case unknow
}

func FXFetchCurrentLanguageCode() -> String {
    let language = Bundle.main.preferredLocalizations.first
    guard let language = language else {
        return "en"
    }
    
    let current = language
    return current
}

func FXCurrentLanguage() -> FXLanguage {
    let language = FXFetchCurrentLanguageCode()
    switch language {
    case "zh-Hans":
        return .zh
    case "zh":
        return .zh
    case "en":
        return .en
    case "fr":
        return .fr
    case "es":
        return .es
    case "it":
        return .it
    case "id":
        return .id
    case "de":
        return .de
    case "ko":
        return .ko
    case "tr":
        return .tr
    case "pt-BR":
        return .pt_BR
    default:
        return .unknow
    }
}


/// 日期本地化
/// - Parameters:
///   - stamp: 时间戳 单位：秒
///   - style: 时间格式
///
///          short:     日期格式：12/13/52
///          medium:    日期格式：Jan 12, 1952
///          long:      日期格式：January 12, 1952
///          full:      日期格式：Tuesday, April 12, 1952 AD
///
/// - Returns: 时间字符串
func FXLocalizationDateWithStyle(stamp: Double, style: DateFormatter.Style = .long) -> String {
    let dateFormatter = DateFormatter()
    dateFormatter.dateStyle = style
    dateFormatter.dateFormat = DateFormatter.dateFormat(fromTemplate: "MMMM-yyyy", options: 0, locale: dateFormatter.locale)
    
    let date = Date(timeIntervalSince1970: stamp)
    let formatterString = dateFormatter.string(from: date)
    return formatterString
}



/// 日期本地化
/// - Parameters:
///   - stamp: 时间戳 单位：秒
///   - template: 时间格式
///
///         "yyyy": 4位数年份，例如 2023
///         "yy":   2位数年份，例如23
///         "MMMM": 完整月份，例如December
///         "MMM":  缩写月份，例如Dec
///         "MM":   2位数月份，例如12
///         "dd":   2位数日期，例如12
///         "EEEE": 完整的星期几，例如Monday
///         "EEE":  缩写星期几，例如Mon
///         "HH":   24小时制，例如12
///         "hh":   12小时制，例如12
///         "mm":   2位数分钟，不足两位加0，例如59
///         "m":    1到2位数分钟，例如，1或者59
///         "zzzz": 完整时区名，例如Pacific Standard Time
///         "zzz":  缩写时区名，例如PST
///         "a":    上午或者下午，例如AM PM
///         "A":    上午或者下午，例如上午 下午
///
/// - Returns: 时间字符串
func FXLocalizationDateWithTemplate(stamp: Double, template: String = "MMMM-yyyy") -> String {
    let dateFormatter = DateFormatter()
    dateFormatter.dateFormat = DateFormatter.dateFormat(fromTemplate: template, options: 0, locale: dateFormatter.locale)
    
    let date = Date(timeIntervalSince1970: stamp)
    let formatterString = dateFormatter.string(from: date)
    return formatterString
}
