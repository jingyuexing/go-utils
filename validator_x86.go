// +build 386
//go:build 386

package utils

func IsIntegerNotMax(val int) bool {
	return val < 0x7FFFFFFF
}
