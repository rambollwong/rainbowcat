package util

import (
	"math/rand"

	"github.com/rambollwong/rainbowcat/types"
)

// SliceContains returns true if an element is present in a collection.
func SliceContains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}
	return false
}

// SliceContainsOneBy returns true if predicate function return true.
func SliceContainsOneBy[T any](collection []T, predicate func(item T) bool) bool {
	for _, item := range collection {
		if predicate(item) {
			return true
		}
	}
	return false
}

// SliceContainsAll returns true if all elements of a subset are contained into a collection or if the subset is empty.
func SliceContainsAll[T comparable](collection []T, subset []T) bool {
	collectionMap := make(map[T]struct{}, len(subset))
	for _, t := range collection {
		collectionMap[t] = struct{}{}
	}
	for _, elem := range subset {
		if _, ok := collectionMap[elem]; !ok {
			return false
		}
	}
	return true
}

// SliceContainsAllBy returns false if predicate function return false.
func SliceContainsAllBy[T any](collection []T, predicate func(item T) bool) bool {
	for _, item := range collection {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// SliceContainsOneOf returns true if at least 1 element of a subset is contained into a collection.
// If the subset is empty SliceContainsOneOf returns false.
func SliceContainsOneOf[T comparable](collection []T, subset []T) bool {
	subMap := make(map[T]struct{}, len(subset))
	for _, t := range subset {
		subMap[t] = struct{}{}
	}
	for _, elem := range collection {
		if _, ok := subMap[elem]; ok {
			return true
		}
	}
	return false
}

// SliceContainsNoneBy returns false if predicate function return true.
func SliceContainsNoneBy[T any](collection []T, predicate func(item T) bool) bool {
	for _, item := range collection {
		if predicate(item) {
			return false
		}
	}
	return true
}

func sliceIntersect[T comparable](list1 []T, list2 []T) []T {
	result := make([]T, 0, len(list1))
	seen := make(map[T]struct{}, len(list1))
	for _, t := range list1 {
		seen[t] = struct{}{}
	}
	for _, t := range list2 {
		if _, ok := seen[t]; ok {
			result = append(result, t)
		}
	}
	return result
}

// SliceIntersect returns the intersection between the collections.
func SliceIntersect[T comparable](list1 []T, listOthers ...[]T) []T {
	if len(listOthers) == 0 {
		return list1
	}
	result := list1
	for _, other := range listOthers {
		result = sliceIntersect(result, other)
	}
	return result
}

// SliceExcludeAll returns slice excluding all given values.
func SliceExcludeAll[T comparable](collection []T, exclude ...T) []T {
	result := make([]T, 0, len(collection))
	excludeMap := make(map[T]struct{}, len(exclude))
	for _, t := range exclude {
		excludeMap[t] = struct{}{}
	}
	for _, t := range collection {
		if _, ok := excludeMap[t]; !ok {
			result = append(result, t)
		}
	}
	return result
}

// SliceExcludeEmpty returns slice excluding empty values.
func SliceExcludeEmpty[T comparable](collection []T) []T {
	result := make([]T, 0, len(collection))
	var empty T
	for _, t := range collection {
		if t != empty {
			result = append(result, t)
		}
	}
	return result
}

// SliceDifference returns the difference between two collections.
// The first value is the collection of element absent of list2.
// The second value is the collection of element absent of list1.
func SliceDifference[T comparable](list1, list2 []T) ([]T, []T) {
	return SliceExcludeAll(list1, list2...), SliceExcludeAll(list2, list1...)
}

// SliceUnion returns all distinct elements from given collections.
// result returns will not change the order of elements relatively.
func SliceUnion[T comparable](lists ...[]T) []T {
	result := make([]T, 0)
	seen := map[T]struct{}{}
	for _, list := range lists {
		for _, e := range list {
			if _, ok := seen[e]; !ok {
				seen[e] = struct{}{}
				result = append(result, e)
			}
		}
	}
	return result
}

// SliceUnionBy returns a duplicate-free version of an array, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is
// invoked for each element in array to generate the criterion by which uniqueness is computed.
func SliceUnionBy[T any, U comparable](iteratee func(index int, item T) U, lists ...[]T) []T {
	result := make([]T, 0)
	seen := map[U]struct{}{}
	for _, list := range lists {
		for i, e := range list {
			u := iteratee(i, e)
			if _, ok := seen[u]; !ok {
				seen[u] = struct{}{}
				result = append(result, e)
			}
		}
	}
	return result
}

// SliceFilter iterates over elements of collection, returning an array of all elements predicate returns truthy for.
func SliceFilter[T any](collection []T, predicate func(index int, item T) bool) []T {
	result := make([]T, 0, len(collection))
	for i, item := range collection {
		if predicate(i, item) {
			result = append(result, item)
		}
	}
	return result
}

// SliceTransformType manipulates a slice and transforms it to a slice of another type.
func SliceTransformType[T any, R any](collection []T, transformer func(index int, item T) R) []R {
	result := make([]R, 0, len(collection))
	for i, item := range collection {
		result = append(result, transformer(i, item))
	}
	return result
}

// SliceFilterTransformType returns a slice which obtained after both filtering
// and transforming using the given callback function.
// The callback function should return two values:
//   - the result of the transforming operation and
//   - whether the result element should be included or not.
func SliceFilterTransformType[T any, R any](collection []T, callback func(index int, item T) (R, bool)) []R {
	result := make([]R, 0, len(collection))
	for i, item := range collection {
		if r, ok := callback(i, item); ok {
			result = append(result, r)
		}
	}
	return result
}

// SliceFlatten returns an array a single level deep.
func SliceFlatten[T any](collection [][]T) []T {
	totalLen := 0
	for i := range collection {
		totalLen += len(collection[i])
	}
	result := make([]T, 0, totalLen)
	for i := range collection {
		result = append(result, collection[i]...)
	}
	return result
}

// SliceFlattenTransformType manipulates a slice and transforms and flattens it to a slice of another type.
// The flatten transformer function can either return a slice or a `nil`, and in the `nil` case
// no value is added to the final slice.
func SliceFlattenTransformType[T any, R any](collection []T, flattenTransformer func(index int, item T) []R) []R {
	result := make([]R, 0, len(collection))
	for i, item := range collection {
		result = append(result, flattenTransformer(i, item)...)
	}
	return result
}

// SliceReduce reduces collection to a value which is the accumulated result of running each element in collection
// through accumulator, where each successive invocation is supplied the return value of the previous.
func SliceReduce[T any, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) R {
	for i, item := range collection {
		initial = accumulator(initial, item, i)
	}
	return initial
}

// SliceReduceRight helper is like Reduce except that it iterates over elements of collection from right to left.
func SliceReduceRight[T any, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) R {
	for i := len(collection) - 1; i >= 0; i-- {
		initial = accumulator(initial, collection[i], i)
	}
	return initial
}

// SliceGroupBy returns an object composed of keys generated
// from the results of running each element of collection through iteratee.
func SliceGroupBy[T any, U comparable](collection []T, iteratee func(item T) U) map[U][]T {
	result := map[U][]T{}
	for _, item := range collection {
		key := iteratee(item)
		result[key] = append(result[key], item)
	}
	return result
}

// SliceOrderedGroupBy returns an array of elements split into groups. The order of grouped values is
// determined by the order they occur in collection. The grouping is generated from the results
// of running each element of collection through iteratee.
func SliceOrderedGroupBy[T any, K comparable](collection []T, iteratee func(item T) K) [][]T {
	result := make([][]T, 0, 1)
	seen := map[K]int{}
	for _, item := range collection {
		key := iteratee(item)
		resultIndex, ok := seen[key]
		if !ok {
			resultIndex = len(result)
			seen[key] = resultIndex
			result = append(result, []T{})
		}
		result[resultIndex] = append(result[resultIndex], item)
	}
	return result
}

// SliceCutChunks returns an array of elements split into groups the length of size. If array can't be split evenly,
// the final chunk will be the remaining elements.
func SliceCutChunks[T any](collection []T, size int) [][]T {
	if size <= 0 {
		panic("Size parameter must be greater than 0")
	}
	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum += 1
	}
	result := make([][]T, 0, chunksNum)
	for i := 0; i < chunksNum; i++ {
		lastIndex := (i + 1) * size
		if lastIndex > len(collection) {
			lastIndex = len(collection)
		}
		result = append(result, collection[i*size:lastIndex])
	}
	return result
}

// SliceInterleaveFlatten round-robin alternating input slices and sequentially appending value at index into result.
func SliceInterleaveFlatten[T any](collections ...[]T) []T {
	if len(collections) == 0 {
		return []T{}
	}
	maxSize := 0
	totalSize := 0
	for _, c := range collections {
		size := len(c)
		totalSize += size
		if size > maxSize {
			maxSize = size
		}
	}
	if maxSize == 0 {
		return []T{}
	}
	result := make([]T, totalSize)
	resultIdx := 0
	for i := 0; i < maxSize; i++ {
		for j := range collections {
			if len(collections[j])-1 < i {
				continue
			}

			result[resultIdx] = collections[j][i]
			resultIdx++
		}
	}
	return result
}

// SliceShuffle returns an array of shuffled values. Uses the Fisher-Yates shuffle algorithm.
func SliceShuffle[T any](collection []T) []T {
	rand.Shuffle(len(collection), func(i, j int) {
		collection[i], collection[j] = collection[j], collection[i]
	})
	return collection
}

// SliceReverse reverses array so that the first element becomes the last,
// the second element becomes the second to last, and so on.
func SliceReverse[T any](collection []T) []T {
	length := len(collection)
	half := length / 2
	for i := 0; i < half; i = i + 1 {
		j := length - 1 - i
		collection[i], collection[j] = collection[j], collection[i]
	}
	return collection
}

// SliceFill fills elements of array with `initial` value.
func SliceFill[T types.Clonable[T]](collection []T, initial T) []T {
	result := make([]T, 0, len(collection))
	for range collection {
		result = append(result, initial.Clone())
	}
	return result
}

// SliceRepeat builds a slice with N copies of initial value.
func SliceRepeat[T types.Clonable[T]](count int, initial T) []T {
	result := make([]T, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, initial.Clone())
	}
	return result
}

// SliceRepeatBy builds a slice with values returned by N calls of callback.
func SliceRepeatBy[T any](count int, predicate func(index int) T) []T {
	result := make([]T, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, predicate(i))
	}
	return result
}

// SliceToMap returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.
func SliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))
	for _, t := range collection {
		k, v := transform(t)
		result[k] = v
	}
	return result
}

// SliceCutLeft drops n elements from the beginning of a slice or array.
// The slice returned is a new slice.
func SliceCutLeft[T any](collection []T, n int) []T {
	if len(collection) <= n {
		return make([]T, 0)
	}
	result := make([]T, 0, len(collection)-n)
	return append(result, collection[n:]...)
}

// SliceCutRight drops n elements from the end of a slice or array.
// The slice returned is a new slice.
func SliceCutRight[T any](collection []T, n int) []T {
	if len(collection) <= n {
		return []T{}
	}
	result := make([]T, 0, len(collection)-n)
	return append(result, collection[:len(collection)-n]...)
}

// SliceCutLeftOn drops elements from the beginning of a slice or array while the predicate returns true.
func SliceCutLeftOn[T any](collection []T, predicate func(item T) bool) []T {
	i := 0
	for ; i < len(collection); i++ {
		if predicate(collection[i]) {
			break
		}
	}
	result := make([]T, 0, len(collection)-i)
	return append(result, collection[i:]...)
}

// SliceCutRightOn drops elements from the end of a slice or array while the predicate returns true.
func SliceCutRightOn[T any](collection []T, predicate func(item T) bool) []T {
	i := len(collection) - 1
	for ; i >= 0; i-- {
		if predicate(collection[i]) {
			break
		}
	}
	result := make([]T, 0, i+1)
	return append(result, collection[:i+1]...)
}

// SliceValueCount counts the number of elements in the collection that compare equal to value.
func SliceValueCount[T comparable](collection []T, value T) (count int) {
	for _, item := range collection {
		if item == value {
			count++
		}
	}
	return count
}

// SliceValueCountBy counts the number of elements in the collection for which predicate is true.
func SliceValueCountBy[T any](collection []T, predicate func(item T) bool) (count int) {
	for _, item := range collection {
		if predicate(item) {
			count++
		}
	}
	return count
}

// SliceValuesCount counts the number of each element in the collection.
func SliceValuesCount[T comparable](collection []T) map[T]int {
	result := make(map[T]int)
	for _, item := range collection {
		result[item]++
	}
	return result
}

// SliceValuesCountBy counts the number of each element return from mapper function.
// Is equivalent to chaining lo.Map and lo.CountValues.
func SliceValuesCountBy[T any, U comparable](collection []T, mapper func(item T) U) map[U]int {
	result := make(map[U]int)
	for _, item := range collection {
		result[mapper(item)]++
	}
	return result
}

// SliceSubset returns a copy of a slice from `offset` up to `length` elements.
// Like `slice[start:start+length]`, but does not panic on overflow.
func SliceSubset[T any](collection []T, offset int, length uint) []T {
	size := len(collection)
	if offset < 0 {
		offset = size + offset
		if offset < 0 {
			offset = 0
		}
	}
	if offset > size {
		return []T{}
	}
	if length > uint(size)-uint(offset) {
		length = uint(size - offset)
	}
	return collection[offset : offset+int(length)]
}

// SliceParagraph returns a copy of a slice from `start` up to, but not including `end`.
// Like `slice[start:end]`, but does not panic on overflow.
func SliceParagraph[T any](collection []T, start int, end int) []T {
	size := len(collection)
	if start >= end {
		return []T{}
	}
	if start > size {
		start = size
	}
	if start < 0 {
		start = 0
	}
	if end > size {
		end = size
	}
	if end < 0 {
		end = 0
	}
	return collection[start:end]
}

// SliceReplace returns a copy of the slice with the first n non-overlapping instances of old replaced by new.
func SliceReplace[T comparable](collection []T, old T, new T, n int) []T {
	result := make([]T, len(collection))
	copy(result, collection)
	for i := range result {
		if result[i] == old && n != 0 {
			result[i] = new
			n--
		}
	}
	return result
}

// SliceReplaceAll returns a copy of the slice with all non-overlapping instances of old replaced by new.
func SliceReplaceAll[T comparable](collection []T, old T, new T) []T {
	return SliceReplace(collection, old, new, -1)
}
