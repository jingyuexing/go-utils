package utils

func IsIntegerNotMax(val int) bool {
	return val < 0x7FFFFFFFFFFFFFFF
}

func IsNonNegative(val int) bool {
	return val >= 0
}

func IsNormarNumber(val int) bool {
	return IsIntegerNotMax(val) && IsNonNegative(val) && IsNotZero(val)
}

func IsNotZero(val int) bool {
	return val != 0
}

func LessThan(val int, less int) bool {
	return val < less
}
func MoreThan(val int, more int) bool {
	return val > more
}

func Enum(val any, enums []any) bool {
	for _, ele := range enums {
		if val == ele {
			return true
		}
	}
	return false
}

func Compose[T any](target T, validators ...(func(T) bool)) bool {
	for _, v := range validators {
		if !v(target) {
			return false
		}
	}
	return true
}
