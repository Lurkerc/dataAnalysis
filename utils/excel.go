package utils

import (
	"strconv"
	"github.com/360EntSecGroup-Skylar/excelize"
)

/*
 * Title : 数组下标转换成excel坐标
 * Author : YC
 * Date : 2018-08-26
 */
func ChangIndexToAxis(intIndexX int, intIndexY int) string {
	var arr = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	intIndexY = intIndexY + 1
	resultY := ""
	for true {
		if intIndexY <= 26 {
			resultY = resultY + arr[intIndexY-1]
			break
		}
		mo := intIndexY % 26
		resultY = arr[mo-1] + resultY
		shang := intIndexY / 26
		if shang <= 26 {
			resultY = arr[shang-1] + resultY
			break
		}
		intIndexY = shang
	}
	return resultY + strconv.Itoa(intIndexX+1)
}

/*
 * Title : 修改excel表格里的值
 * Author : YC
 * Date : 2018-08-26
 */
func ModifyExcelCellByAxis(xlsx *excelize.File, sheet string, axis string, value interface{}) int {
	xlsx.SetCellValue(sheet, axis, value)
	return 0
}
