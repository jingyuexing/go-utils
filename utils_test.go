package utils_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	utils "github.com/jingyuexing/go-utils"
)

func TestFormat(t *testing.T) {
	result := utils.Template("{a} {b}", map[string]any{
		"a": "hello",
		"b": "world",
	}, "")
	if result != "hello world" {
		t.Error("not expected result", result)
	} else {
		fmt.Printf(result + "\n")
	}

	result2 := utils.Template("? ?", map[string]any{
		"a": 1,
		"b": 2,
	}, "?")
	if result2 != "1 2" {
		t.Error(fmt.Sprintf("Expect %s, but got %s", "1 2", result2))
	}

	result3 := utils.Template("/a /b", map[string]any{
		"a": "hello",
		"b": "world",
	}, "/")

	if result3 != "hello world" {
		t.Error(fmt.Sprintf("Expect %s, but got %s", "hello world", result3))
	}

	result4 := utils.Template("【a】 【b】", map[string]any{
		"a": "hello",
		"b": "world",
	}, "【】")

	if result4 != "hello world" {
		t.Error(fmt.Sprintf("Expect %s, but got %s", "hello world", result4))
	}

	result5 := utils.Template("《a》 《b》", map[string]any{
		"a": "hello",
		"b": "world",
	}, "《》")

	if result5 != "hello world" {
		t.Error(fmt.Sprintf("Expect %s, but got %s", "hello world", result5))
	}

	result6 := utils.Template("！a！ ！b！", map[string]any{
		"a": "hello",
		"b": "world",
	}, "！！")

	if result6 != "hello world" {
		t.Error(fmt.Sprintf("Expect %s, but got %s", "hello world", result6))
	}

	result7 := utils.Template("👉a👈 👉b👈", map[string]any{
		"a": "hello",
		"b": "world",
	}, "👉👈")

	if result7 != "hello world" {
		t.Error(fmt.Sprintf("Expect %s, but got %s", "hello world", result7))
	}

	result8 := utils.Template("[a] [b]", map[string]any{
		"a": "hello",
		"b": "world",
	}, "[]")

	if result8 != "hello world" {
		t.Error(fmt.Sprintf("Expect %s, but got %s", "hello world", result8))
	}
}

func TestCookieParse(t *testing.T) {
	cookie := new(utils.Cookie)
	cookie.NewCookie("a=b&c=23&k=66", "&", "=")
	cookie.PutOne("s", "v")
	fmt.Printf("%s\n", cookie.ToString())
}

func TestDebounce(t *testing.T) {

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
	toggle1 := utils.UseToggle("on", "off")
	toggle1.Switch()
	if toggle1.Value() != "off" {
		t.Error("not pass")
	}
	if toggle1.Switch() != "on" {
		t.Error("not pass")
	}

	toggle2 := utils.UseToggle(true, false)
	if toggle2.Value() != true {
		t.Error("not pass")
	}

	if toggle2.Switch() != false {
		t.Error("not pass")
	}

	toggle3 := utils.UseToggle(1, 2, 3, 4, 5, 6, 7, 8, 9)
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
	}, 2)
}

func TestReduce(t *testing.T) {
	sum := utils.Reduce([]int{1, 2, 3, 5, 6}, func(accumulator, currentValue int) int {
		return accumulator + currentValue
	}, 0)
	if sum != 17 {
		t.Error("not pass")
	}
}

func TestOption(t *testing.T) {
	type S struct{}
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

func TestEmit(t *testing.T) {
	event := utils.NewEventEmit()
	count := 1
	event.On("plus", func(args ...any) {
		count += (args[0]).(int)
	})

	for i := 1; i <= 3; i++ {
		event.Emit("plus", i)
	}

	if count != 7 {
		t.Error("not pass")
	}

}

func TestAccessNested(t *testing.T) {
	type S struct {
		A struct {
			B int
		}
	}
	k := &S{
		A: struct{ B int }{
			B: 33,
		},
	}
	value := utils.AccessNested(k, "A.B", ".")
	if value.(int) != 33 {
		t.Error("not pass")
	}
}

func TestNestedAttr(t *testing.T) {
	type C struct {
		Value string
	}
	type B struct {
		C *C
	}
	type A struct {
		B B
	}

	testAData := &A{
		B{
			&C{
				Value: "Haaaa",
			},
		},
	}
	var finalValue any
	err := utils.NestedObject(testAData, "B.C.Value", func(target reflect.Value, key string) {
		finalValue = target.Interface()
		fmt.Println("Value:", finalValue, "key", key)
	})
	if err != nil {
		t.Fatal("has wrong")
	}

	if finalValue.(string) != "Haaaa" {
		t.Error(fmt.Sprintf("Expect: %s, but got %s", "Haaaa", finalValue.(string)))
	}

	testData2 := map[string]any{
		"A": map[string]any{
			"b": map[string]any{
				"C": map[string]any{
					"Value": "Haaaa",
				},
			},
		},
	}
	var finalValue2 any
	utils.NestedObject(testData2, "A.b.C.Value", func(target reflect.Value, key string) {
		finalValue2 = target.Interface()
		fmt.Println("Value:", finalValue2, "key", key)

	})

	if finalValue2.(string) != "Haaaa" {
		t.Error(fmt.Sprintf("Expect: %s, but got %s", "Haaaa", finalValue2.(string)))
	}
}

func TestDateTimeFormat(t *testing.T) {
	datetime := time.Unix(1706572292, 0)
	year := utils.DateTimeFormat(datetime, "YYYY")
	if year != "2024" {
		t.Error("format year is wrong")
	}

	result1 := utils.DateTimeFormat(datetime, "YY/MM/dd")
	if result1 != "24/01/30" {
		t.Error("format date wrong")
	}

	result2 := utils.DateTimeFormat(datetime, "MM/dd/YYYY HH:mm:ss")

	if result2 != "01/30/2024 07:51:32" {
		t.Error("format datetime wrong")
	}

	result3 := utils.DateTimeFormat(datetime, "YYYY年M月d日 H时m分s秒")
	if result3 != "2024年1月30日 7时51分32秒" {
		t.Error("format has wrong")
	}
}

func TestDateTime(t *testing.T) {
	datetime := utils.NewDateTime()
	datetime = datetime.SetTime(1706572292,0)

    if datetime.Day != 30 {
        t.Error("the field Day is wrong")
    }

    if datetime.Year != 2024 {
        t.Error("the field Day is wrong")
    }

    if datetime.Add(2,"year").Year != 2026 {
        t.Error("the add 2 year has wrong")
    }

    if datetime.Add(1,"week").Day != 6 {
        t.Error("the add 1 week has wrong")
    }
    fmt.Println(datetime.String())


    if datetime.String() != "2024-01-30T07:51:32.5132Z" {
		t.Error("format time has wrong")
	}

    if datetime.Add(100,"D").String() != "2024-05-09T07:51:32.5132Z" {
        t.Error("add 100 days has wrong")
    }

    if datetime.CurrentYearDays() != 366 {
        t.Error("caculate has wrong")
    }

}
