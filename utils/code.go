package utils

import "github.com/axgle/mahonia"

// @Title 字符串编码转换
// @Params srcString string 文本
// @Params srcCode string 原编码格式
// @Params tagCode string 目标编码格式
// ConvertToString("文本", "gbk", "utf-8")
func ConvertToString(srcString, srcCode, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(srcString)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, data, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(data)
	return result
}
