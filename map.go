package goutil

// Keys returns all keys from a map.
func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Values returns all values from a map.
func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// ContainsKey returns true if the key exists in the map.
func ContainsKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// Pick returns a new map containing only the specified keys.
func Pick[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}
	return result
}

// Omit returns a new map excluding the specified keys.
func Omit[K comparable, V any](m map[K]V, keys []K) map[K]V {
	excludeSet := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		excludeSet[k] = struct{}{}
	}
	result := make(map[K]V, len(m))
	for k, v := range m {
		if _, skip := excludeSet[k]; !skip {
			result[k] = v
		}
	}
	return result
}

// Merge merges multiple maps into one. Later maps override earlier ones.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// MapValues transforms map values using the iteratee function.
func MapValues[K comparable, V any, U any](m map[K]V, iteratee func(V) U) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = iteratee(v)
	}
	return result
}

// MapKeys transforms map keys using the iteratee function.
// If duplicate keys result from transformation, the last value wins.
func MapKeys[K comparable, V any, J comparable](m map[K]V, iteratee func(K) J) map[J]V {
	result := make(map[J]V, len(m))
	for k, v := range m {
		result[iteratee(k)] = v
	}
	return result
}

// FilterMap returns a new map containing only entries that satisfy the predicate.
func FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// ForEachMap iterates over a map and invokes the function for each entry.
func ForEachMap[K comparable, V any](m map[K]V, iteratee func(K, V)) {
	for k, v := range m {
		iteratee(k, v)
	}
}

// Invert swaps keys and values. Values must be comparable.
// If duplicate values exist in the source map, one of them wins (non-deterministic).
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// Entries converts a map to a slice of key-value pairs.
func Entries[K comparable, V any](m map[K]V) []Entry[K, V] {
	result := make([]Entry[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, Entry[K, V]{Key: k, Value: v})
	}
	return result
}

// FromEntries converts a slice of key-value pairs to a map.
func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V {
	result := make(map[K]V, len(entries))
	for _, e := range entries {
		result[e.Key] = e.Value
	}
	return result
}

// Entry represents a key-value pair.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// FindKey returns the first key whose value equals the given value.
// Returns zero key and false if not found.
func FindKey[K comparable, V comparable](m map[K]V, value V) (K, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// FindKeyBy returns the first key where the predicate returns true.
// Returns zero key and false if not found.
func FindKeyBy[K comparable, V any](m map[K]V, predicate func(K, V) bool) (K, bool) {
	for k, v := range m {
		if predicate(k, v) {
			return k, true
		}
	}
	var zero K
	return zero, false
}
