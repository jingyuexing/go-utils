package utils_test

import (
	"fmt"
	"testing"

	utils "github.com/jingyuexing/go-utils"
)

func TestPipelCallback(t *testing.T){
    addOne := func(x int) int {
        return x + 1
    }

    double := func(x int) int {
        return x * 2
    }

    square := func(x int) int {
        return x * x
    }
    composeFunc := utils.PipeCallback(addOne,double,square)

    result := composeFunc(55)
    if result != 12544 {
        t.Error("no pass")
    }else{
        fmt.Printf("%d",result)
    }
}
