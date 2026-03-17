package goutil

import (
	"cmp"
	"time"
)

// Index iterates from 0 to limit-1 and returns the first index
// where predicate returns true. Returns -1 if no match is found.
func Index(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// Contains returns true if the slice contains the specified item.
func Contains[T comparable](collection []T, element T) bool {
	for _, v := range collection {
		if v == element {
			return true
		}
	}
	return false
}

// ContainsBy returns true if any element matches the predicate.
func ContainsBy[T any](collection []T, predicate func(T) bool) bool {
	for _, v := range collection {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Filter returns elements that satisfy the predicate.
func Filter[T any](collection []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range collection {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reject returns elements that do NOT satisfy the predicate (opposite of Filter).
func Reject[T any](collection []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range collection {
		if !predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map transforms each element using the iteratee function.
func Map[T any, U any](collection []T, iteratee func(T) U) []U {
	result := make([]U, len(collection))
	for i, v := range collection {
		result[i] = iteratee(v)
	}
	return result
}

// Reduce reduces a collection to a single value using an accumulator function.
func Reduce[T any, U any](collection []T, accumulator func(U, T) U, initial U) U {
	result := initial
	for _, v := range collection {
		result = accumulator(result, v)
	}
	return result
}

// ForEach iterates over a collection and invokes the function for each element.
func ForEach[T any](collection []T, iteratee func(T, int)) {
	for i, v := range collection {
		iteratee(v, i)
	}
}

// Find returns the first element matching the predicate and true.
// Returns zero value and false if no match is found.
func Find[T any](collection []T, predicate func(T) bool) (T, bool) {
	for _, v := range collection {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FindIndexOf returns the element, its index, and true for the first match.
// Returns zero value, -1, and false if not found.
func FindIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool) {
	for i, v := range collection {
		if predicate(v) {
			return v, i, true
		}
	}
	var zero T
	return zero, -1, false
}

// IndexOf returns the index of the first occurrence of element, or -1 if not found.
func IndexOf[T comparable](collection []T, element T) int {
	for i, v := range collection {
		if v == element {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of the last occurrence of element, or -1 if not found.
func LastIndexOf[T comparable](collection []T, element T) int {
	for i := len(collection) - 1; i >= 0; i-- {
		if collection[i] == element {
			return i
		}
	}
	return -1
}

// FindLastIndexOf returns the element, its index, and true for the last match.
// Returns zero value, -1, and false if not found.
func FindLastIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool) {
	for i := len(collection) - 1; i >= 0; i-- {
		if predicate(collection[i]) {
			return collection[i], i, true
		}
	}
	var zero T
	return zero, -1, false
}

// FindOrElse returns the first element matching the predicate,
// or fallback if no match is found.
func FindOrElse[T any](collection []T, fallback T, predicate func(T) bool) T {
	for _, v := range collection {
		if predicate(v) {
			return v
		}
	}
	return fallback
}

// FindUniques returns elements that appear exactly once in the collection.
func FindUniques[T comparable](collection []T) []T {
	counts := make(map[T]int, len(collection))
	for _, v := range collection {
		counts[v]++
	}
	result := make([]T, 0)
	for _, v := range collection {
		if counts[v] == 1 {
			result = append(result, v)
		}
	}
	return result
}

// FindUniquesBy returns elements whose key (from keyFn) appears exactly once.
func FindUniquesBy[T any, K comparable](collection []T, keyFn func(T) K) []T {
	counts := make(map[K]int, len(collection))
	for _, v := range collection {
		counts[keyFn(v)]++
	}
	result := make([]T, 0)
	for _, v := range collection {
		if counts[keyFn(v)] == 1 {
			result = append(result, v)
		}
	}
	return result
}

// FindDuplicates returns one copy of each element that appears more than once.
func FindDuplicates[T comparable](collection []T) []T {
	counts := make(map[T]int, len(collection))
	for _, v := range collection {
		counts[v]++
	}
	seen := make(map[T]struct{})
	result := make([]T, 0)
	for _, v := range collection {
		if counts[v] > 1 {
			if _, exists := seen[v]; !exists {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// FindDuplicatesBy returns one copy of each element whose key appears more than once.
func FindDuplicatesBy[T any, K comparable](collection []T, keyFn func(T) K) []T {
	counts := make(map[K]int, len(collection))
	for _, v := range collection {
		counts[keyFn(v)]++
	}
	seen := make(map[K]struct{})
	result := make([]T, 0)
	for _, v := range collection {
		key := keyFn(v)
		if counts[key] > 1 {
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// Unique returns a new slice with duplicates removed, preserving order.
func Unique[T comparable](collection []T) []T {
	seen := make(map[T]struct{}, len(collection))
	result := make([]T, 0, len(collection))
	for _, v := range collection {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// UniqueBy returns a unique slice using a key function to determine equality.
func UniqueBy[T any, K comparable](collection []T, keyFn func(T) K) []T {
	seen := make(map[K]struct{}, len(collection))
	result := make([]T, 0, len(collection))
	for _, v := range collection {
		key := keyFn(v)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// GroupBy groups elements by a key returned from the keyFn.
func GroupBy[T any, K comparable](collection []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range collection {
		key := keyFn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// KeyBy transforms a slice into a map using the given key function.
// If duplicate keys exist, the last value wins.
func KeyBy[T any, K comparable](collection []T, keyFn func(T) K) map[K]T {
	result := make(map[K]T, len(collection))
	for _, v := range collection {
		result[keyFn(v)] = v
	}
	return result
}

// Partition splits a slice into two: elements matching the predicate,
// and elements that don't.
func Partition[T any](collection []T, predicate func(T) bool) (matched []T, unmatched []T) {
	matched = make([]T, 0)
	unmatched = make([]T, 0)
	for _, v := range collection {
		if predicate(v) {
			matched = append(matched, v)
		} else {
			unmatched = append(unmatched, v)
		}
	}
	return matched, unmatched
}

// Chunk splits a slice into chunks of the specified size.
func Chunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	chunks := make([][]T, 0, (len(collection)+size-1)/size)
	for i := 0; i < len(collection); i += size {
		end := i + size
		if end > len(collection) {
			end = len(collection)
		}
		chunks = append(chunks, collection[i:end])
	}
	return chunks
}

// Flatten merges a slice of slices into a single slice.
func Flatten[T any](collections [][]T) []T {
	total := 0
	for _, s := range collections {
		total += len(s)
	}
	result := make([]T, 0, total)
	for _, s := range collections {
		result = append(result, s...)
	}
	return result
}

// FlatMap applies iteratee to each element and flattens the result.
func FlatMap[T any, U any](collection []T, iteratee func(T) []U) []U {
	result := make([]U, 0)
	for _, v := range collection {
		result = append(result, iteratee(v)...)
	}
	return result
}

// Compact removes zero-value elements from a slice.
func Compact[T comparable](collection []T) []T {
	var zero T
	result := make([]T, 0, len(collection))
	for _, v := range collection {
		if v != zero {
			result = append(result, v)
		}
	}
	return result
}

// Reverse returns a new slice with elements in reverse order.
func Reverse[T any](collection []T) []T {
	result := make([]T, len(collection))
	for i, v := range collection {
		result[len(collection)-1-i] = v
	}
	return result
}

// Without returns a slice excluding the specified values.
func Without[T comparable](collection []T, exclude ...T) []T {
	excludeSet := make(map[T]struct{}, len(exclude))
	for _, v := range exclude {
		excludeSet[v] = struct{}{}
	}
	result := make([]T, 0, len(collection))
	for _, v := range collection {
		if _, skip := excludeSet[v]; !skip {
			result = append(result, v)
		}
	}
	return result
}

// Intersect returns elements that exist in both slices.
func Intersect[T comparable](list1 []T, list2 []T) []T {
	set := make(map[T]struct{}, len(list2))
	for _, v := range list2 {
		set[v] = struct{}{}
	}
	result := make([]T, 0)
	for _, v := range list1 {
		if _, exists := set[v]; exists {
			result = append(result, v)
		}
	}
	return result
}

// Difference returns elements in list1 that are NOT in list2.
func Difference[T comparable](list1 []T, list2 []T) []T {
	set := make(map[T]struct{}, len(list2))
	for _, v := range list2 {
		set[v] = struct{}{}
	}
	result := make([]T, 0)
	for _, v := range list1 {
		if _, exists := set[v]; !exists {
			result = append(result, v)
		}
	}
	return result
}

// Union returns all unique elements from multiple slices.
func Union[T comparable](lists ...[]T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)
	for _, list := range lists {
		for _, v := range list {
			if _, exists := seen[v]; !exists {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// Count returns the number of elements equal to the given value.
func Count[T comparable](collection []T, value T) int {
	count := 0
	for _, v := range collection {
		if v == value {
			count++
		}
	}
	return count
}

// CountBy returns the number of elements matching the predicate.
func CountBy[T any](collection []T, predicate func(T) bool) int {
	count := 0
	for _, v := range collection {
		if predicate(v) {
			count++
		}
	}
	return count
}

// Min returns the minimum value. Returns zero value if empty.
func Min[T cmp.Ordered](collection []T) T {
	if len(collection) == 0 {
		var zero T
		return zero
	}
	minVal := collection[0]
	for _, v := range collection[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

// MinIndex returns the minimum value and its index. Returns (zero, -1) if empty.
func MinIndex[T cmp.Ordered](collection []T) (T, int) {
	if len(collection) == 0 {
		var zero T
		return zero, -1
	}
	minVal := collection[0]
	minIdx := 0
	for i, v := range collection[1:] {
		if v < minVal {
			minVal = v
			minIdx = i + 1
		}
	}
	return minVal, minIdx
}

// MinBy returns the minimum element using a comparison function.
// comparison(a, b) should return true when a < b.
func MinBy[T any](collection []T, comparison func(a, b T) bool) T {
	if len(collection) == 0 {
		var zero T
		return zero
	}
	minVal := collection[0]
	for _, v := range collection[1:] {
		if comparison(v, minVal) {
			minVal = v
		}
	}
	return minVal
}

// MinIndexBy returns the minimum element and its index using a comparison function.
// comparison(a, b) should return true when a < b.
func MinIndexBy[T any](collection []T, comparison func(a, b T) bool) (T, int) {
	if len(collection) == 0 {
		var zero T
		return zero, -1
	}
	minVal := collection[0]
	minIdx := 0
	for i, v := range collection[1:] {
		if comparison(v, minVal) {
			minVal = v
			minIdx = i + 1
		}
	}
	return minVal, minIdx
}

// Max returns the maximum value. Returns zero value if empty.
func Max[T cmp.Ordered](collection []T) T {
	if len(collection) == 0 {
		var zero T
		return zero
	}
	maxVal := collection[0]
	for _, v := range collection[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// MaxIndex returns the maximum value and its index. Returns (zero, -1) if empty.
func MaxIndex[T cmp.Ordered](collection []T) (T, int) {
	if len(collection) == 0 {
		var zero T
		return zero, -1
	}
	maxVal := collection[0]
	maxIdx := 0
	for i, v := range collection[1:] {
		if v > maxVal {
			maxVal = v
			maxIdx = i + 1
		}
	}
	return maxVal, maxIdx
}

// MaxBy returns the maximum element using a comparison function.
// comparison(a, b) should return true when a > b.
func MaxBy[T any](collection []T, comparison func(a, b T) bool) T {
	if len(collection) == 0 {
		var zero T
		return zero
	}
	maxVal := collection[0]
	for _, v := range collection[1:] {
		if comparison(v, maxVal) {
			maxVal = v
		}
	}
	return maxVal
}

// MaxIndexBy returns the maximum element and its index using a comparison function.
// comparison(a, b) should return true when a > b.
func MaxIndexBy[T any](collection []T, comparison func(a, b T) bool) (T, int) {
	if len(collection) == 0 {
		var zero T
		return zero, -1
	}
	maxVal := collection[0]
	maxIdx := 0
	for i, v := range collection[1:] {
		if comparison(v, maxVal) {
			maxVal = v
			maxIdx = i + 1
		}
	}
	return maxVal, maxIdx
}

// Latest returns the most recent time.Time from the collection.
// Returns zero time if empty.
func Latest(collection []time.Time) time.Time {
	if len(collection) == 0 {
		return time.Time{}
	}
	latest := collection[0]
	for _, v := range collection[1:] {
		if v.After(latest) {
			latest = v
		}
	}
	return latest
}

// LatestBy returns the element with the most recent time, extracted by timeFn.
// Returns zero value if empty.
func LatestBy[T any](collection []T, timeFn func(T) time.Time) T {
	if len(collection) == 0 {
		var zero T
		return zero
	}
	latest := collection[0]
	latestTime := timeFn(latest)
	for _, v := range collection[1:] {
		t := timeFn(v)
		if t.After(latestTime) {
			latest = v
			latestTime = t
		}
	}
	return latest
}

// Every returns true if all elements match the predicate.
func Every[T any](collection []T, predicate func(T) bool) bool {
	for _, v := range collection {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Some returns true if at least one element matches the predicate.
func Some[T any](collection []T, predicate func(T) bool) bool {
	for _, v := range collection {
		if predicate(v) {
			return true
		}
	}
	return false
}

// None returns true if no element matches the predicate.
func None[T any](collection []T, predicate func(T) bool) bool {
	for _, v := range collection {
		if predicate(v) {
			return false
		}
	}
	return true
}

// First returns the first element and true, or zero value and false if empty.
func First[T any](collection []T) (T, bool) {
	if len(collection) == 0 {
		var zero T
		return zero, false
	}
	return collection[0], true
}

// FirstOrEmpty returns the first element, or zero value if the slice is empty.
func FirstOrEmpty[T any](collection []T) T {
	if len(collection) == 0 {
		var zero T
		return zero
	}
	return collection[0]
}

// Last returns the last element and true, or zero value and false if empty.
func Last[T any](collection []T) (T, bool) {
	if len(collection) == 0 {
		var zero T
		return zero, false
	}
	return collection[len(collection)-1], true
}

// LastOrEmpty returns the last element, or zero value if the slice is empty.
func LastOrEmpty[T any](collection []T) T {
	if len(collection) == 0 {
		var zero T
		return zero
	}
	return collection[len(collection)-1]
}
