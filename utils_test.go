package utils_test

import (
	"fmt"
	"testing"

	utils "github.com/jingyuexing/go-utils"
)

func TestFormat(t *testing.T) {
	result := utils.Template("{a} {b}", map[string]any{
		"a": "hello",
		"b": "world",
	})
	if result != "hello world" {
		t.Error("not expected result", result)
	} else {
		fmt.Printf(result + "\n")
	}
}

func TestCookieParse(t *testing.T) {
	cookie := new(utils.Cookie)
	cookie.NewCookie("a=b&c=23&k=66", "&", "=")
	cookie.PutOne("s", "v")
	fmt.Printf("%s\n", cookie.ToString())
}

func TestgetPathValue(t *testing.T) {
	result := utils.GetPathValue("/user/:id", "/user/12")
	if result["id"] != "12" {
		t.Error("Not Pass")
	}
	result2 := utils.GetPathValue("/user/{id}", "/user/456")
	if result2["id"] != "456" {
		t.Error("NOT PASS")
	}
}

func TestMap2Struct(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	user := &User{}
	utils.Map2Struct(map[string]any{
		"name": "V",
		"age":  10,
	}, user)
	if user.Name != "V" {
		t.Error("Not Pass")
	}
}

func TestToggle(t *testing.T) {
    toggle1 := utils.UseToggle("on","off")
    toggle1.Switch()
    if toggle1.Value() != "off" {
        t.Error("not pass")
    }
    if toggle1.Switch() != "on" {
        t.Error("not pass")
    }

    toggle2 := utils.UseToggle(true,false)
    if toggle2.Value() != true {
        t.Error("not pass")
    }

    if toggle2.Switch() != false {
        t.Error("not pass")
    }

    toggle3 := utils.UseToggle(1,2,3,4,5,6,7,8,9)
    if toggle3.Value() != 1 {
        t.Error("not pass")
    }
    toggle3.Switch()
    toggle3.Switch()
    toggle3.Switch()
    toggle3.Switch()
    if toggle3.Value() != 5 {
        t.Error("not pass")
    }

}


func TestTimes(t *testing.T) {
     utils.Times(func(t ...string) string {
        return "hello world"
     },2)
}

func TestReduce(t *testing.T) {
    sum := utils.Reduce([]int{1,2,3,5,6},func(accumulator, currentValue int) int {
        return accumulator + currentValue
    },0)
    if sum != 17 {
        t.Error("not pass")
    }
}

func TestOption(t *testing.T) {
    type S struct {}
    v := &S{}
    v = nil
    raw := utils.Some(v)

    if !raw.IsNone() {
        t.Error("not pass")
    }
    if raw.UnwrapOr(12).(int) != 12 {
        t.Error("not pass")
    }

}

func TestEmit(t *testing.T){
    event := utils.NewEventEmit()
    count := 1
    event.On("plus",func(args ...any) {
        count+= (args[0]).(int)
    })

    for i := 1; i <= 3; i++ {
        event.Emit("plus",i)
    }

    if count != 7 {
        t.Error("not pass")
    }

}
