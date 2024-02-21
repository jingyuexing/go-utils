package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ToString(val any) string {
	valstr := ""
	switch v := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		valstr = fmt.Sprintf("%d", v)
	case float32, float64:
		valstr = fmt.Sprintf("%f", v)
	case bool:
		valstr = fmt.Sprintf("%v", v)
	case string:
		valstr = fmt.Sprintf("%s", v)
	default:
		valstr = fmt.Sprintf("%v", v)
	}
	return valstr
}

func Template(source string, data map[string]any, placeholder string) string {
	tempSyntax := "{}"
	if placeholder != "" {
		tempSyntax = placeholder
	}
	key := ""
	sourceCopy := &source
	for k, val := range data {
		valStr := ""
		valStr = ToString(val)
		switch tempSyntax {
		case "{}":
			key = strings.Join([]string{"{", k, "}"}, "")
		case "[]":
			key = strings.Join([]string{"[", k, "]"}, "")
		case "()":
			key = strings.Join([]string{"(", k, ")"}, "")
		case ":", "@", "#", "$":
			key = strings.Join([]string{tempSyntax, k}, "")
		case "?":
			key = tempSyntax
		default:
			if (len(tempSyntax) & 1) == 1 {
				key = strings.Join([]string{tempSyntax, k}, "")
				break
			}
			if (len(tempSyntax) & 1) == 0 {
				halfLength := len(tempSyntax) / 2
				key = strings.Join([]string{tempSyntax[:halfLength], k, tempSyntax[halfLength:]}, "")
			}

		}
		*sourceCopy = strings.Replace(*sourceCopy, key, valStr, 1)
	}
	return *sourceCopy
}

func InRange(val int, min int, max int) bool {
	return val >= min && val <= max
}

func IsUpper(char rune) bool {
	return InRange(int(char), int('A'), int('Z'))
}

func IsNumberic(char rune) bool {
	return InRange(int(char), int('0'), int('9'))
}
func IsLower(char rune) bool {
	return InRange(int(char), int('a'), int('z'))
}

func IsAlphabet(char rune) bool {
	return IsLower(char) || IsUpper(char)
}

func IsAlphanum(char rune) bool {
	return IsAlphabet(char) || IsNumberic(char)
}

func getValidRef(refname string, symbol rune) string {
	word := []string{}
	for _, char := range refname {
		if char != symbol && int(char) != ' ' {
			if IsAlphanum(char) {
				word = append(word, string(char))
			} else if int(char) > 0xff {
				word = append(word, string(char))
			}
		}
	}
	return strings.Join(word, "")
}

func toComment(val string, commentSymbol []string) string{
    if len(commentSymbol) & 1 == 0 {
        return strings.Join([]string{commentSymbol[0],val,commentSymbol[1]},"")
    }
    return commentSymbol[0] + val;
}

func deepReplace(target map[string]string, symbol rune, text string) string {
	clone := strings.Clone(text)
	list := strings.Split(clone, " ")
	for i := 0; i < len(list); i++ {
		token := list[i]
		if !strings.Contains(clone, string(symbol)) {
			return clone
		}
		if strings.Contains(token, string(symbol)) {
			refrence := strings.Split(token, string(symbol))
			for _, refName := range refrence {
				if refName == "" {
					continue
				}
				validRef := getValidRef(refName, rune(symbol))
				refString, ok := target[validRef]
				if !ok {
					refString = toComment(validRef,[]string{"/*","*/"})
				}
				values := deepReplace(target, symbol, refString)
				clone = strings.Replace(clone, string(symbol)+refName, values, 1)
			}
		}
	}
	return clone
}

func ReferenceString(target map[string]string, symbol rune) func(string) string {
	if symbol == 0 {
		symbol = '@'
	}
	replace := func(text string) string {
		return deepReplace(target, symbol, text)
	}
	getItem := func(key string) string {
		if val, ok := target[key]; ok {
			return replace(val)
		}
		return ""
	}
	return getItem
}

func FindVariableNames(text string, delimiter string) []string {
	leftSyombol := ""
	rightSyombol := ""
	variables := []string{}
	maxLength := len(text)
	if maxLength&0 == 0 {
		half := len(delimiter) / 2
		leftSyombol = delimiter[:half]
		rightSyombol = delimiter[half:]
	} else {
		leftSyombol = delimiter
	}
	idx := 0
	for idx < maxLength {
		nameCache := make([]string, 0)
		if string(text[idx]) == leftSyombol {
			n := idx + 1
			for n < maxLength {
				if rightSyombol != "" && string(text[n]) == rightSyombol {
					break
				} else if string(text[n]) == " " {
					break
				}
				nameCache = append(nameCache, string(text[n]))
				n++
			}
			variables = append(variables, strings.Join(nameCache, ""))
			//
			nameCache = make([]string, 0)
			idx = n
		}
		idx++
	}
	return variables
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

func CreateNestedObject(keys string, value interface{}, delimiter string) map[string]interface{} {
	if delimiter == "" {
		delimiter = "."
	}
	keyList := strings.Split(keys, delimiter)
	result := make(map[string]interface{})
	current := reflect.ValueOf(result)

	for i, key := range keyList {
		if i == len(keyList)-1 {
			current.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value).Elem())
		} else {
			if current.MapIndex(reflect.ValueOf(key)).IsValid() {
				current = current.MapIndex(reflect.ValueOf(key))
			} else {
				newMap := make(map[string]interface{})
				current.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(newMap))
				current = reflect.ValueOf(newMap)
			}
		}
	}

	return result
}

func AccessNested(data any, path string, delimiter string) any {
	keys := strings.Split(path, delimiter)
	value := reflect.ValueOf(data)

	if value.Kind() == reflect.Pointer {
		value = reflect.Indirect(value)
	}

	for _, key := range keys {
		if !value.IsValid() {
			return nil
		}

		switch value.Kind() {
		case reflect.Pointer:
			value = reflect.Indirect(value)
		case reflect.Map:
			mapValue := value.MapIndex(reflect.ValueOf(key))
			if !mapValue.IsValid() {
				return nil
			}
			value = mapValue
		case reflect.Struct:
			field := value.FieldByName(key)
			if !field.IsValid() {
				return nil
			}
			value = field
		default:
			return nil
		}
	}

	if value.IsValid() && value.CanInterface() {
		return value.Interface()
	}
	return nil
}

func NestedObject[T any](target T, path string, cb func(target reflect.Value, key string)) error {
	keys := strings.Split(path, ".")
	v := reflect.ValueOf(target)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for _, key := range keys {
		switch v.Kind() {
		case reflect.Map:
			mapValue := v.MapIndex(reflect.ValueOf(key))

			if mapValue.Kind() == reflect.Interface ||
				mapValue.Kind() == reflect.Pointer {
				mapValue = mapValue.Elem()
			}
			if mapValue.IsValid() {
				cb(mapValue, key)
				if mapValue.Kind() == reflect.Map || mapValue.Kind() == reflect.Pointer {
					v = mapValue
				}
			} else {
				return errors.New("unreachable access path")
			}
		case reflect.Struct:
			field := v.FieldByName(key)
			if field.Kind() == reflect.Pointer ||
				field.Kind() == reflect.Interface ||
				field.Kind() == reflect.Map {
				field = field.Elem()
			}
			if field.IsValid() {
				v = field
				cb(v, key)
			} else {
				return errors.New("unreachable access path")
			}
		case reflect.Interface:
			if v.Elem().Kind() == reflect.Map || v.Elem().Kind() == reflect.Pointer {
				v = v.Elem()
			} else {
				return errors.New("unreachable access path")
			}
		case reflect.Pointer:
			if v.Elem().IsValid() {
				v = v.Elem()
			} else {
				return errors.New("unreachable access path")
			}
		default:
			return errors.New("invalid target type")
		}
	}

	return nil
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

func TimeDuration(duration string) (time.Time, error) {
	const (
		SECOND        uint64 = 1
		MINUTE_SECOND uint64 = 60 * SECOND
		HOUR_SECOND   uint64 = 60 * MINUTE_SECOND
		DAY_SECOND    uint64 = 24 * HOUR_SECOND
		WEEK_SECOND   uint64 = 7 * DAY_SECOND
		MONTH_SECOND  uint64 = 30 * DAY_SECOND
		YEAR_SECOND   uint64 = 12 * MONTH_SECOND
	)
	uints := map[string]uint64{
		"y": YEAR_SECOND,
		"m": MONTH_SECOND,
		"w": WEEK_SECOND,
		"d": DAY_SECOND,
		"h": HOUR_SECOND,
		"s": SECOND,
	}
	splitDuration := strings.Split(duration, " ")
	durationSecond := uint64(0)
	for _, val := range splitDuration {
		unit := val[len(val)-1:]
		num, _ := strconv.Atoi(val[:len(val)-1])
		if unitSeconds, ok := uints[unit]; ok {
			durationSecond += uint64(num) * unitSeconds
		} else {
			return time.Time{}, fmt.Errorf("haven't this unit")
		}
	}
	return time.Now().Add(time.Duration(durationSecond) * time.Second).UTC(), nil
}
