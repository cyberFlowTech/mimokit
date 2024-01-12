package lan

import "golang.org/x/text/language"

// SupportLan 语言包文件名与这里一致
var SupportLan = map[string]language.Tag{
	"en":    language.English,
	"zh_CN": language.Chinese,
	"zh_TW": language.TraditionalChinese,
	"ja_JP": language.Japanese,
	"ko_KR": language.Korean,
	"vi_VN": language.Vietnamese,
}
