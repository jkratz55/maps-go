package maps

import (
	"fmt"
)

// ConflictResolver is a function type that is invoked when Merge is merging maps that
// contains the same key, hence a conflict. The role of ConflictResolver is to handle
// the conflict and return the resolved value.
type ConflictResolver[V any] func(left, right V) V

// OverwriteResolver returns a ConflictResolver that always overrides the existing value.
func OverwriteResolver[V any]() ConflictResolver[V] {
	return func(left, right V) V {
		return right
	}
}

// NopResolver returns a ConflictResolver that always keeps the existing value.
func NopResolver[V any]() ConflictResolver[V] {
	return func(left, right V) V {
		return left
	}
}

// Merge merges multiple maps into a single new map. If a key exists in multiple maps the
// ConflictResolver function is called to resolve the conflict. The value returns by the
// ConflictResolver is the value set in the new merged map.
func Merge[M ~map[K]V, K comparable, V any](fn ConflictResolver[V], src ...M) map[K]V {
	merged := make(map[K]V)
	for _, m := range src {
		for k, v := range m {
			if existing, ok := merged[k]; ok {
				newVal := fn(existing, v)
				merged[k] = newVal
			} else {
				merged[k] = v
			}
		}
	}
	return merged
}

// Keys returns all the keys in the provided map.
//
// The keys will be in an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values return all the values in the provided map.
//
// The values will be in an indeterminate order.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// GetOrDefault returns the value for a given key in the provided map, or returns
// the default value if the key doesn't exist.
func GetOrDefault[M ~map[K]V, K comparable, V any](m M, key K, defaultVal V) V {
	if val, ok := m[key]; ok {
		return val
	}
	return defaultVal
}

// GetOrPanic returns the value for a given key in the provided map or panics if
// the key doesn't exist in the map.
func GetOrPanic[M ~map[K]V, K comparable, V any](m M, key K) V {
	if val, ok := m[key]; ok {
		return val
	}
	panic(fmt.Errorf("key %v doesn't exist in map", key))
}

// SetIfPresent sets the value for the key in the map only if the key already exist
// in the map. If the value is set SetIfPresent returns true, otherwise returns false.
func SetIfPresent[M ~map[K]V, K comparable, V any](m M, key K, val V) bool {
	if _, ok := m[key]; ok {
		m[key] = val
		return true
	}
	return false
}

// SetIfAbsent sets the value for the key in the map only if the key does not already
// exist in the map. If the value is set SetIfAbsent returns true, otherwise returns
// false.
func SetIfAbsent[M ~map[K]V, K comparable, V any](m M, key K, val V) bool {
	if _, ok := m[key]; !ok {
		m[key] = val
		return true
	}
	return false
}

// Clear removes all entries from the map.
func Clear[M ~map[K]V, K comparable, V any](m M) {
	for key := range m {
		delete(m, key)
	}
}

// Clone clones a map adding all the entries of the map into a new map.
//
// Note: If V is a pointer type or contains types backed by a pointer or maps,
// slices, channels, functions, etc. this will not be a deep copy.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	if m == nil {
		return nil
	}
	newMap := make(M, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// Copy copies all the entries from the source map into the destination map. If
// the key already exist in the dest map its value will be overwritten.
func Copy[M ~map[K]V, K comparable, V any](src, dst M) {
	for k, v := range src {
		dst[k] = v
	}
}

// Equal compares two maps and returns a boolean value indicating if they are equal.
func Equal[M ~map[K]V, K, V comparable](m1, m2 M) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

// Entry is a data structure representing a single entry in a map.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Entries returns all entries in the given map as a slice of Entry.
//
// The results will be in an indeterminate order.
func Entries[M ~map[K]V, K comparable, V any](m M) []Entry[K, V] {
	res := make([]Entry[K, V], 0, len(m))
	for k, v := range m {
		res = append(res, Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return res
}

// EntryMapper is a function type that maps an entry from a map to a new map, possibly of different types.
type EntryMapper[K1, K2 comparable, V1, V2 any] func(key K1, val V1) (K2, V2)

// MapEntries manipulates a maps entries and transforms it to another map.
func MapEntries[M1 ~map[K1]V1, M2 ~map[K2]V2, K1, K2 comparable, V1, V2 any](in M1, mapper EntryMapper[K1, K2, V1, V2]) M2 {
	res := make(M2, len(in))
	for k1, v1 := range in {
		k2, v2 := mapper(k1, v1)
		res[k2] = v2
	}
	return res
}

// Predicate represents a predicate (boolean-value function).
type Predicate[K comparable, V any] func(key K, val V) bool

// Filter iterates through the entries of the map and tests if they satisfy the predicate.
// Entries that satisfy the predicate added to a newly returned map, effectively filtering
// the map entries.
func Filter[M ~map[K]V, K comparable, V any](m M, fn Predicate[K, V]) M {
	res := make(M)
	for k, v := range m {
		if fn(k, v) {
			res[k] = v
		}
	}
	return res
}
