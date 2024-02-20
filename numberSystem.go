package utils

import (
	"math"
	"strings"
)

func ToChineseNumber(number int64, base int, upper bool) string {
	digits := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	units := []string{"", "十", "百", "千", "万", "亿"}
    trans := [...]string{"廿","卅"}
	if upper {
		digits = []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
		units = []string{"", "拾", "佰", "仟", "万", "亿"}
	}

	result := ""

	quotient := int64(math.Abs(float64(number)))
	if base > len(digits) {
		return ""
	}

	unitIndex := 0
	for quotient > 0 {
		remainder := quotient % int64(base)
		if remainder > 0 || unitIndex == 0 || (unitIndex == 1 && quotient%10 == 0) {
			if unitIndex == 1 && (remainder > 1 && remainder <= 3) {
				result = trans[remainder - 2] + result
			} else {
				if remainder == 0 && unitIndex > 0 {
					result = digits[remainder] + result
				} else if unitIndex/4 > 1 {
					unitIndex = unitIndex % 4
					result = digits[remainder] + (units[unitIndex % 4] + units[unitIndex]) + result
				} else {
					result = digits[remainder] + units[unitIndex] + result

				}
			}
		}
		quotient = int64(math.Floor(float64(quotient / int64(base))))
		unitIndex++
	}

	if number < 0 {
		return "负" + result
	} else if result == "" {
		return digits[0] // 处理数字为零的情况
	} else {
		return removeTrailingZero(result)
	}
}

func removeTrailingZero(str string) string {
	for strings.HasSuffix(str,"零") {
        str = strings.Replace(str,"零","",1)
    }
	return str
}
