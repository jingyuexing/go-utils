package utils

import "strings"


const (
    hextable = "0123456789abcdef"
)

type Buffer []byte

func NewBuffer() Buffer {
    return make([]byte,0)
}

func (b Buffer) Length() int{
    return len(b) // the unint is byte
}

func (b Buffer) String() string {
    str := make([]byte,b.Length() << 1)
    j := 0
    for _,item := range b {
        str[j] = hextable[item>>4]
        str[j+1] = hextable[item&0x0f]
        j += 2
    }
    return string(str)
}

func (b Buffer) Load(data string) Buffer{
    for i:= 0;i < len(data);i+=2 {
        cache := 0x00
        cache = strings.Index(hextable,string(data[i])) << 4
        cache = cache | strings.Index(hextable,string(data[i+1]))
        b = append(b,  byte(cache))
    }
    return b
}
