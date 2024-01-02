## go-utils

This is a method or function that I often use


## usage

```go
func Template(source string, data map[string]string, placeholder string) string
```
- This function can perform template string substitution

[see test case](utils_test.go#L10-L82)

```go
func GetPathValue(raw string, realPath string) map[string]string
```
- This string will directly get the value corresponding to the template [see test case](utils_test.go#L91-L100)

```go
func (c *Cookie) NewCookie(cookie string, delimiter string, joiner string) *Cookie
```
- This can parse the Cookie string or parse the Query string, you can specify separators and connectors for parsing [see test case](utils_test.go#L84-L89)


```go
func Map2Struct(source map[string]any, bindingTarget any)
```
- You can bind the value of the map to the structure, provided that the fields of the structure and the fields of the map are consistent in lowercase [see test case](utils_test.go#L102-L105)


