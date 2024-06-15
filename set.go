package utils

type Set[T comparable] map[T]bool

func (s Set[T]) Values() []T {
	values := make([]T, 0)
	for k, _ := range s {
		values = append(values, k)
	}
	return values
}

func (s Set[T]) Has(value T) bool {
    _,ok := s[value]
    return ok
}

func (s Set[T]) Add(value T) bool {
    s[value] = true
    return s[value]
}

func NewSet[T comparable](value ...T) Set[T] {
    set := make(Set[T],0)
    for _,val := range value {
        set.Add(val)
    }
    return set
}



// Union 返回两个集合的并集
func (s Set[T]) Union(other Set[T]) Set[T] {
    unionSet := NewSet[T]()
    for k := range s {
        unionSet.Add(k)
    }
    for k := range other {
        unionSet.Add(k)
    }
    return unionSet
}

func (s Set[T]) Intersection(other Set[T]) Set[T] {
    intersectionSet := NewSet[T]()
    for k := range s {
        if other.Has(k) {
            intersectionSet.Add(k)
        }
    }
    return intersectionSet
}

func (s Set[T]) Difference(other Set[T]) Set[T] {
    differenceSet := NewSet[T]()
    for k := range s {
        if !other.Has(k) {
            differenceSet.Add(k)
        }
    }
    return differenceSet
}
