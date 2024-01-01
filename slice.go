package utils


func Chunk[T any](slice []T, size int) [][]T {
    var chunks [][]T
    for i := 0; i < len(slice); i += size {
        end := i + size
        if end > len(slice) {
            end = len(slice)
        }
        chunks = append(chunks, slice[i:end])
    }
    return chunks
}


func Reduce[T any](slice []T, reducer func(T, T) T, initialValue T) T {
    result := initialValue

    for _, val := range slice {
        result = reducer(result, val)
    }

    return result
}

func Map[T any](slice []T, mapper func(T) T) []T {
    result := make([]T, len(slice))

    for i, val := range slice {
        result[i] = mapper(val)
    }

    return result
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T

    for _, val := range slice {
        if predicate(val) {
            result = append(result, val)
        }
    }

    return result
}

func PadEnd[T any](slice []T, targetLength int, padValue T) []T {
    currentLength := len(slice)
    paddingLength := targetLength - currentLength
    if paddingLength <= 0 {
        return slice
    }
    newSlice := make([]T, 0,paddingLength)
    for i := 0; i < paddingLength; i++ {
        newSlice = append(newSlice, padValue)
    }
    newSlice = append(slice, newSlice...)

    return newSlice
}

func PadStart[T any](slice []T, targetLength int, padValue T) []T {
    currentLength := len(slice)
    paddingLength := targetLength - currentLength
    if paddingLength <= 0 {
        return slice
    }

    newSlice := make([]T, 0, paddingLength)
    for i := 0; i < paddingLength; i++ {
        newSlice = append(newSlice, padValue)
    }

    newSlice = append(newSlice, slice...)

    return newSlice
}

func GroupBy[T any, K comparable](slice []T, getKey func(T) K) map[K][]T {
    groups := make(map[K][]T)

    for _, item := range slice {
        key := getKey(item)
        groups[key] = append(groups[key], item)
    }

    return groups
}

func ForEach[T any](slice []T,callback func(value T,key int)){
    for k, v := range slice {
        callback(v,k)
    }
}
