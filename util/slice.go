package util

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
