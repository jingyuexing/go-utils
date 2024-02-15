package utils

import "reflect"

// omit specify field in target struct
//
//  type People struct {
//     Name    string
//     Age     int
//     Address string
// }
// p := &People{
//     Name:    "william",
//     Age:     20,
//     Address: "BC",
// }
//  Omit(p,[]string{"Name"})
//
func Omit[T any](target T, fields []string) map[string]any {
	mapping := map[string]bool{}
	for _, value := range fields {
		mapping[value] = true
	}
	result := StructFilter(target, func(field string, val reflect.Value) bool {
		_, ok := mapping[field]
		return !ok
	})
	return result
}
// pick specify field in target struct
//
//  type People struct {
//     Name    string
//     Age     int
//     Address string
// }
// p := &People{
//     Name:    "william",
//     Age:     20,
//     Address: "BC",
// }
//  Pick(p,[]string{"Name"})
//
func Pick[T any](target T, fields []string) map[string]any {
	mapping := map[string]bool{}
	for _, value := range fields {
		mapping[value] = true
	}
	result := StructFilter(target, func(field string, val reflect.Value) bool {
		_, ok := mapping[field]
		return ok
	})
	return result
}

func StructFilter[T any](target T, callback func(field string, val reflect.Value) bool) map[string]any {
	reflectValue := reflect.ValueOf(target)
	reflectType := reflect.TypeOf(target)
	if reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
        reflectType = reflectType.Elem()
	}
	result := make(map[string]any)
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectType.Field(i)
		if callback(field.Name, reflectValue.FieldByName(field.Name)) {
			result[field.Name] = reflectValue.FieldByName(field.Name).Interface()
		}
	}
	return result
}
