//go:build amd64 || arm64
// +build amd64 arm64

package utils

func IsIntegerNotMax(val int) bool {
	return val < 0x7FFFFFFFFFFFFFFF
}
