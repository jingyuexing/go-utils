package utils

import (
	"errors"
	"fmt"
)

// PipeCallback executes a series of functions where the output of each function is passed as input to the next function.
// Returns a function that can be invoked to run the entire execution process.
func PipeCallback[T any](fns ...func(args T) T) func(args T) T {
	return func(args T) T {
		result := args
		for _, fn := range fns {
			result = fn(result)
		}
		return result
	}
}

func Times[T any](cb func(args ...T) T, x int) func(...T) T {
	var value T
	times := 0
	return func(args ...T) T {
		if times < x {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = errors.New(fmt.Sprint(r))
					}
					fmt.Println("Error occurred:", err)
				}
			}()
			value = cb(args...)
			times++
		}

		return value
	}
}

func Compose[T any](funcs ...func(args ...T) T) func(args ...T) T {
    callback := func(a func(args ...T) T, b func(args ...T) T) func(args ...T) T {
        return func(args ...T) T {
            return a(b(args...))
        }
    }
    startFunc := funcs[0]
    for i := 1; i < len(funcs); i++ {
        startFunc = callback(startFunc, funcs[i])
    }
    return startFunc
}
