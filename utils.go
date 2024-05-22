package mytrade

import "strings"

// 计算小数点后有效位数的函数
func countDecimalPlaces(str string) int {
	// 去除尾部的零
	str = strings.TrimRight(str, "0")
	// 分割字符串以获取小数部分
	parts := strings.Split(str, ".")
	// 如果有小数部分，则返回其长度
	if len(parts) == 2 {
		return len(parts[1])
	}
	// 如果没有小数部分，则返回0
	return 0
}
