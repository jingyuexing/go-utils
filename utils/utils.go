package utils

import (
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
    query map[string]string
    delimiter string
    joiner string
}


func (c *Cookie)NewCookie(cookie string,delimiter string, joiner string) *Cookie{
    c.delimiter = delimiter
    c.joiner = joiner
    var cookiesList = strings.Split(cookie,c.delimiter)
    c.query = make(map[string]string)
    for _,item :=  range cookiesList{
        keyAndValue := strings.Split(item,c.joiner)
        c.query[keyAndValue[0]] = keyAndValue[1]
    }
    return c
}

func (c *Cookie)PutOne(key string,val string) *Cookie{
    c.query[key] = val
    return c
}

func (c *Cookie)GetAll()map[string]string{
    return c.query
}

func (c *Cookie)ToString()string{
    var cookieSplit []string
    for key,val := range c.query{
        cookieSplit = append(cookieSplit, strings.Join([]string{key,val},c.joiner))
    }
    return strings.Join(cookieSplit,c.delimiter)
}
