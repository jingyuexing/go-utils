package utils

type Options[T any] struct {
    values T
    none   bool
}

type OptionInterface[T any] interface {
    IsSome() bool
    IsNone() bool
    Unwrap() T
    Expect(msg string) T
    UnwrapOr(defaultValue T) T
    IsSomeAnd(fn func(args ...interface{}) bool) bool
    Inspect() Options[T]
}

func Some[T any](value T) Options[T] {
    return Options[T]{values: value, none: false}
}

func None[T any]() Options[T] {
    return Options[T]{none: true}
}

func (o Options[T]) IsSome() bool {
    return !o.none
}

func (o Options[T]) IsNone() bool {
    return o.none
}

func (o Options[T]) Unwrap() T {
    if o.none {
        panic("Unwrap: Option is None")
    }
    return o.values
}

func (o Options[T]) Expect(msg string) T {
    if o.none {
        panic("Expect: " + msg)
    }
    return o.values
}

func (o Options[T]) UnwrapOr(defaultValue T) T {
    if o.none {
        return defaultValue
    }
    return o.values
}

func (o Options[T]) IsSomeAnd(fn func(args ...interface{}) bool) bool {
    if o.none {
        return false
    }
    return fn(o.values)
}

func (o Options[T]) Inspect() Options[T] {
    return o
}

func Option[T any](value T) Options[T] {
    if(&value == nil){
        return None[T]()
    }else{
        return Some[T](value)
    }
}
type People struct {
    Age int
    Name string
    Address string
}

func main(){
    var age *int = nil;
    ok := Option(age)
    vals := *ok.Expect("this is nil pointer")
    somePeople := Option(&People{
        Age: 20,
        Name: "Gog",
        Address: "Mars",
    })
    somePeople.Unwrap()
}
