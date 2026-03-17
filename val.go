package goutil

import "reflect"

// IsNil returns true if the value is nil.
// Handles the Go interface nil gotcha — a typed nil pointer
// wrapped in an interface will correctly return true.
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// IsNotNil returns true if the value is not nil.
func IsNotNil(v any) bool {
	return !IsNil(v)
}

// Empty returns the zero value of type T.
func Empty[T any]() T {
	var zero T
	return zero
}

// IsEmpty returns true if the value equals the zero value of its type.
func IsEmpty[T comparable](v T) bool {
	var zero T
	return v == zero
}

// IsNotEmpty returns true if the value is not the zero value of its type.
func IsNotEmpty[T comparable](v T) bool {
	var zero T
	return v != zero
}
