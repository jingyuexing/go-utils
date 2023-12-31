package utils
type Toggler[T any] interface {
    Switch() T
    Value() T
}

type toggleValue[T any] struct {
    values     []T
    currentIdx int
}

func (t *toggleValue[T]) Switch() T {
    t.currentIdx = (t.currentIdx + 1) % len(t.values)
    return t.Value()
}

func (t *toggleValue[T]) Value() T {
    return t.values[t.currentIdx]
}

func UseToggle[T any](value ...T) Toggler[T] {
    return &toggleValue[T]{values: value, currentIdx: 0}
}
