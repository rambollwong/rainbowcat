package util

import "github.com/rambollwong/rainbowcat/types"

// MapKeys creates an array of the map keys.
func MapKeys[K comparable, V any](in map[K]V) []K {
	result := make([]K, 0, len(in))
	for k := range in {
		result = append(result, k)
	}
	return result
}

// MapValues creates an array of the map values.
func MapValues[K comparable, V any](in map[K]V) []V {
	result := make([]V, 0, len(in))
	for _, v := range in {
		result = append(result, v)
	}
	return result
}

// MapValueOr returns the value of the given key or the fallback value if the key is not present.
func MapValueOr[K comparable, V any](in map[K]V, key K, fallback V) V {
	if v, ok := in[key]; ok {
		return v
	}
	return fallback
}

// MapFilter returns same map type filtered by given predicate.
func MapFilter[K comparable, V any](in map[K]V, predicate func(key K, value V) bool) map[K]V {
	r := map[K]V{}
	for k, v := range in {
		if predicate(k, v) {
			r[k] = v
		}
	}
	return r
}

// MapFilterByKeys returns same map type filtered by given keys.
func MapFilterByKeys[K comparable, V any](in map[K]V, keys []K) map[K]V {
	r := map[K]V{}
	keysMap := SliceToMap(keys, func(item K) (K, struct{}) {
		return item, struct{}{}
	})
	for k, v := range in {
		if _, ok := keysMap[k]; ok {
			r[k] = v
		}
	}
	return r
}

// MapFilterByValues returns same map type filtered by given values.
func MapFilterByValues[K comparable, V comparable](in map[K]V, values []V) map[K]V {
	r := map[K]V{}
	valuesMap := SliceToMap(values, func(item V) (V, struct{}) {
		return item, struct{}{}
	})
	for k, v := range in {
		if _, ok := valuesMap[v]; ok {
			r[k] = v
		}
	}
	return r
}

// MapExcludeByKeys returns a map excluding all given keys.
func MapExcludeByKeys[K comparable, V any](in map[K]V, keys []K) map[K]V {
	r := map[K]V{}
	keysMap := SliceToMap(keys, func(item K) (K, struct{}) {
		return item, struct{}{}
	})
	for k, v := range in {
		if _, ok := keysMap[k]; !ok {
			r[k] = v
		}
	}
	return r
}

// MapExcludeByValues returns a map excluding all given values.
func MapExcludeByValues[K comparable, V comparable](in map[K]V, values []V) map[K]V {
	r := map[K]V{}
	valuesMap := SliceToMap(values, func(item V) (V, struct{}) {
		return item, struct{}{}
	})
	for k, v := range in {
		if _, ok := valuesMap[v]; !ok {
			r[k] = v
		}
	}
	return r
}

// MapEntries transforms a map into array of key/value pairs.
func MapEntries[K comparable, V any](in map[K]V) []types.Entry[K, V] {
	entries := make([]types.Entry[K, V], 0, len(in))
	for k, v := range in {
		entries = append(entries, types.Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return entries
}

// MapFromEntries transforms an array of key/value pairs into a map.
func MapFromEntries[K comparable, V any](entries []types.Entry[K, V]) map[K]V {
	out := make(map[K]V, len(entries))
	for _, v := range entries {
		out[v.Key] = v.Value
	}
	return out
}

// MapInvert creates a map composed of the inverted keys and values. If map
// contains duplicate values, subsequent values overwrite property assignments
// of previous values.
func MapInvert[K comparable, V comparable](in map[K]V) map[V]K {
	out := make(map[V]K, len(in))
	for k, v := range in {
		out[v] = k
	}
	return out
}

// MapAssign merges multiple maps from left to right.
func MapAssign[K comparable, V any](maps ...map[K]V) map[K]V {
	out := map[K]V{}
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// MapTransformKeys manipulates a map keys and transforms it to a map of another type.
func MapTransformKeys[K comparable, V any, R comparable](in map[K]V, iteratee func(value V, key K) R) map[R]V {
	result := make(map[R]V, len(in))
	for k, v := range in {
		result[iteratee(v, k)] = v
	}
	return result
}

// MapTransformValues manipulates a map values and transforms it to a map of another type.
func MapTransformValues[K comparable, V any, R any](in map[K]V, iteratee func(value V, key K) R) map[K]R {
	result := make(map[K]R, len(in))
	for k, v := range in {
		result[k] = iteratee(v, k)
	}
	return result
}

// MapTransformKeyValues manipulates a map entries and transforms it to a map of another type.
func MapTransformKeyValues[K1 comparable, V1 any, K2 comparable, V2 any](
	in map[K1]V1,
	iteratee func(key K1, value V1) (K2, V2),
) map[K2]V2 {
	result := make(map[K2]V2, len(in))
	for k1, v1 := range in {
		k2, v2 := iteratee(k1, v1)
		result[k2] = v2
	}
	return result
}

// MapToSlice transforms a map into a slice based on specific iteratee.
func MapToSlice[K comparable, V any, R any](in map[K]V, iteratee func(key K, value V) R) []R {
	result := make([]R, 0, len(in))

	for k, v := range in {
		result = append(result, iteratee(k, v))
	}

	return result
}
