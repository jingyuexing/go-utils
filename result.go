package utils

import (
	"errors"
	"fmt"
)

type Result[T any] struct {
	value T     // 成功值
	err   error // 错误信息
}

func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}
func (r *Result[T]) IsOk() bool {
	return r.err == nil
}

func (r *Result[T]) Value() T {
	return r.value
}

func (r *Result[T]) Err() error {
	return r.err
}

func (r *Result[T]) IsErr() bool {
	return r.err != nil
}

func (r *Result[T]) Unwrap() T {
	if r.IsErr() {
		panic(r.err)
	}
	return r.value
}

func (r *Result[T]) UnwrapOr(defaultValue T) T {
	if r.IsErr() {
		return defaultValue
	}
	return r.value
}

func (r *Result[T]) TryCatch(callback func(T) error) Result[T] {
	if r.err != nil {
		return *r
	}

	// Wrapper function to capture panic
	var runCallback = func(value T) (result Result[T]) {
		defer func() {
			if err := recover(); err != nil {
				var e error
				switch err := err.(type) {
				case error:
					e = err
				case string:
					e = errors.New(err)
				default:
					e = fmt.Errorf("unknown panic: %v", err)
				}
				result = Err[T](e)
			}
		}()

		err := callback(value)
		if err != nil {
			return Err[T](err)
		}
		return Ok(value)
	}

	return runCallback(r.value)
}

func NewResult[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// This function takes a Result[T] and a mapping function f of type func(T) U
// and returns a new Result[U] containing the transformed value
// func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
//   if !r.ok {
//     return Err[U](r.err)
//   }
//   return Ok[U](f(r.val))
// }

// func (r Result[T]) AndThen(f func(T) Result[U]) Result[U] {
//     if !r.ok {
//         return Err[U](r.err)
//     }
//     return f(r.val)
// }
