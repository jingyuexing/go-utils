package utils_test

import (
	"testing"

	"jingyuexing.com/utils"
)

func TeTestCompose(t *testing.T) {
	if !utils.Compose(12, utils.IsNotZero, utils.IsNonNegative) {
		t.Error("NOT PASS")
	}
}
