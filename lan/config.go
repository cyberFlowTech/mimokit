package lan

import (
	"golang.org/x/text/language"
)

// SupportLan 语言包文件名与这里一致
var SupportLan = map[string]language.Tag{
	"en":    language.English,            // 英文
	"zh_CN": language.Chinese,            // 简体中文
	"zh_TW": language.TraditionalChinese, // 繁体中文
	"ja_JP": language.Japanese,           // 日语
	"ko_KR": language.Korean,             // 朝鲜语
	"vi_VN": language.Vietnamese,         // 越南语
	"mn_MN": language.Mongolian,          // 蒙古语
	"pt_PT": language.Portuguese,         // 葡萄牙语
	"es_ES": language.Spanish,            // 西班牙语
	"th_TH": language.Thai,               // 西班牙语
	"id_ID": language.Indonesian,         // 印尼语
}
