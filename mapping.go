package utils

import (
	"reflect"
	"strings"
	"unicode"
)

// KeyFunc defines a function type for generating map keys from struct fields
type KeyFunc func(field reflect.StructField) string

// omit specify field in target struct
//
//	 type People struct {
//	    Name    string
//	    Age     int
//	    Address string
//	}
//
//	p := &People{
//	    Name:    "william",
//	    Age:     20,
//	    Address: "BC",
//	}
//
//	Omit(p,"Name")
func Omit[T any](target T, fields ...string) map[string]any {
	mapping := map[string]bool{}
	for _, value := range fields {
		mapping[SnakeCase(value)] = true
		mapping[value] = true
		mapping[value] = true
	}
	result := StructFilter(target, func(field string, val reflect.Value) bool {
		_, ok := mapping[field]
		return !ok
	}, nil)
	return result
}

// pick specify field in target struct
//
//	 type People struct {
//	    Name    string
//	    Age     int
//	    Address string
//	}
//
//	p := &People{
//	    Name:    "william",
//	    Age:     20,
//	    Address: "BC",
//	}
//
//	Pick(p,[]string{"Name"})
func Pick[T any](target T, fields ...string) map[string]any {
	mapping := map[string]bool{}
	for _, value := range fields {
		mapping[value] = true
	}
	result := StructFilter(target, func(field string, val reflect.Value) bool {
		_, ok := mapping[field]
		return ok
	}, nil)
	return result
}

func jsonKeyFunc(field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		return jsonTag
	}
	return ""
}

// SnakeCase 将 CamelCase 转换为 snake_case，并处理缩写词
func SnakeCase(s string) string {
	var result []rune
	n := len(s)
	allUpper := true

	for _, r := range s {
		if unicode.IsLower(r) {
			allUpper = false
			break
		}
	}

	if allUpper {
		return strings.ToLower(s)
	}

	for i := 0; i < n; i++ {
		r := rune(s[i])

		if unicode.IsUpper(r) {
			if i > 0 && !unicode.IsUpper(rune(s[i-1])) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// GetFieldName 根据字段的 JSON 标签或字段名返回正确的字段名称
func GetFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return SnakeCase(field.Name)
	}
	jsonName := strings.Split(jsonTag, ",")[0] // 获取 JSON 标签的第一个部分
	if jsonName != "" && jsonName != "-" {
		return jsonName
	}
	return SnakeCase(field.Name)
}

func structFilterRecursive(
	value reflect.Value,
	t reflect.Type,
	callback func(field string, val reflect.Value) bool,
	keyFunc KeyFunc,
	result map[string]any,
) {
	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)
		fieldValue := value.Field(i)

		key := keyFunc(field)

		if field.Anonymous {
			// If the field is an embedded struct, recurse into it
			structFilterRecursive(fieldValue, fieldValue.Type(), callback, keyFunc, result)
		} else if callback(key, fieldValue) && key != "" {
			result[key] = fieldValue.Interface()
		}
	}
}

func StructFilter[T any](
	target T,
	callback func(field string, val reflect.Value) bool,
	key KeyFunc,
) map[string]any {
	reflectValue := reflect.ValueOf(target)
	reflectType := reflect.TypeOf(target)
	if reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
		reflectType = reflectType.Elem()
	}
	GenKey := key
	if key == nil {
		GenKey = GetFieldName
	}
	result := make(map[string]any)
	structFilterRecursive(reflectValue, reflectType, callback, GenKey, result)
	return result
}

func resolvePointer(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Pointer {
		return reflect.Indirect(value)
	}
	return value
}

func GetFields[T any](target T, callback func(reflect.Value)) {
	reflectValue := reflect.ValueOf(target)
	reflectValue = resolvePointer(reflectValue)
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		fieldType := reflectValue.Type().Field(i)
		if fieldType.Anonymous {
			embeddedValue := resolvePointer(field)
			if embeddedValue.Kind() == reflect.Struct {
				field := embeddedValue.FieldByName(field.Type().Name())
				callback(field)
			}
		} else {
			callback(field)
		}
	}
}

// 泛型函数 ArrayFilter，接收一个目标数组 target 和回调函数 cb
func ArrayFilter[T any](target []T, cb func(T) []map[string]any) []map[string]any {
	var result []map[string]any // 存储过滤后的结果

	for _, item := range target {
		// 调用回调函数 cb，对每个元素进行处理
		filtered := cb(item)

		// 如果 cb 函数返回非空的结果，则将其添加到最终结果中
		if len(filtered) > 0 {
			result = append(result, filtered...)
		}
	}

	return result
}
