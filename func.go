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

func Times(cb func(...any) any, x int) func(...any) any {
	value := any(nil)
	times := 0
	return func(args ...any) any {
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
