package utils_test

import (
	"testing"

	utils "github.com/jingyuexing/go-utils"
)

func TestCompose(t *testing.T) {
	if !utils.ValidateCompose(12, utils.IsNotZero, utils.IsNonNegative) {
		t.Error("NOT PASS")
	}
}
