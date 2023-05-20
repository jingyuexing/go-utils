package utils

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
)

func Template(source string, data map[string]string) string {
	sourceCopy := &source
	for k, v := range data {
		*sourceCopy = strings.Replace(*sourceCopy, strings.Join([]string{"{", k, "}"}, ""), v, 1)
	}
	return *sourceCopy
}

type Cookie struct {
	query     map[string]string
	delimiter string
	joiner    string
}

func pathParse(target string) []string {
	PathSlice := strings.Split(target, "/")
	obj := make([]string, 0)
	for i := 0; i < len(PathSlice); i++ {
		key := PathSlice[i]
		if key != "" {
			if strings.HasPrefix(key, ":") {
				if strings.HasSuffix(key, ":") {
					obj = append(obj, key[1:len(key)-1])
				} else {
					obj = append(obj, key[1:])
				}
			} else if strings.HasPrefix(key, "{") && strings.HasSuffix(key, "}") {
				obj = append(obj, key[1:len(key)-1])

			} else {
				obj = append(obj, key)
			}
		}
	}
	return obj
}

func GetPathValue(raw string, realPath string) map[string]string {
	rawObj := pathParse(raw)
	realObj := pathParse(realPath)
	valObj := make(map[string]string)
	for idx, pathName := range rawObj {
		if rawObj[idx] != realObj[idx] {
			valObj[pathName] = realObj[idx]
		}
	}
	return valObj
}

func (c *Cookie) NewCookie(cookie string, delimiter string, joiner string) *Cookie {
	c.delimiter = delimiter
	c.joiner = joiner
	var cookiesList = strings.Split(cookie, c.delimiter)
	c.query = make(map[string]string)
	for _, item := range cookiesList {
		keyAndValue := strings.Split(item, c.joiner)
		c.query[keyAndValue[0]] = keyAndValue[1]
	}
	return c
}

func (c *Cookie) PutOne(key string, val string) *Cookie {
	c.query[key] = val
	return c
}

func (c *Cookie) GetAll() map[string]string {
	return c.query
}

func (c *Cookie) ToString() string {
	var cookieSplit []string
	for key, val := range c.query {
		cookieSplit = append(cookieSplit, strings.Join([]string{key, val}, c.joiner))
	}
	return strings.Join(cookieSplit, c.delimiter)
}

func LoadConfig[T any](filename string) *T {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("connfig file not found or not exites current directory", err)
	}
	var result *T
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal("json parse error", err)
	}
	return result
}

func ToLowerCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func Map2Struct(source map[string]any, bindingTarget any) {
	typeEle := reflect.TypeOf(bindingTarget).Elem()
	ValEle := reflect.ValueOf(bindingTarget).Elem()
	if ValEle.Kind() == reflect.Struct && ValEle.CanSet() {
		for i := 0; i < typeEle.NumField(); i++ {
			field := typeEle.Field(i)
			value, ok := source[ToLowerCamelCase(field.Name)]
			if ok && field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
				ptr := ValEle.FieldByName(field.Name)
				if ptr.IsNil() {
					newStructPtr := reflect.New(field.Type.Elem())
					ptr.Set(newStructPtr)
				}
				Map2Struct(source, ptr.Interface())
			} else {
				ValEle.FieldByName(field.Name).Set(reflect.ValueOf(value))
			}
		}
	}
}
