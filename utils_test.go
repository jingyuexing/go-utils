package utils_test

import (
	"fmt"
	"testing"

	utils "jingyuexing.com/jingyuexing/go-utils"
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
