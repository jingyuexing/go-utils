package utils_test

import (
    "fmt"
    "testing"

    utils "jingyuexing.com/utils"
)

func TestFormat(t *testing.T) {
    result := utils.Template("{a} {b}", map[string]string{
        "a": "hello",
        "b": "world",
    })
    if result != "hello world" {
        t.Error("not expected result", result)
    } else {
        fmt.Printf(result + "\n")
    }
}

func TestCookieParse(t *testing.T){
    cookie := new(utils.Cookie)
    cookie.NewCookie("a=b&c=23&k=66","&","=")
    cookie.PutOne("s","v")
    fmt.Printf("%s\n",cookie.ToString())
}
