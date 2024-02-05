package util

func GetDuplicates(list []string) []string {
	elementCount := CountGrouped(list)

	// get list of names with count != 1
	var duplicates []string
	for name, count := range elementCount {
		if count != 1 {
			duplicates = append(duplicates, name)
		}
	}
	return duplicates
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

// CountGrouped counts the occurrences of equal elements (groups).
// You can use this to determine if there are duplicates or if an element exists.
func CountGrouped[T comparable](elements []T) map[T]int {
	counts := map[T]int{}
	for _, element := range elements {
		counts[element] += 1
	}
	return counts
}

// Any determines if any element in the given list matches the given predicate
func Any[T any](list []T, predicate func(T) bool) bool {
	for _, t := range list {
		if predicate(t) {
			return true
		}
	}
	return false
}

// GroupBy groups elements of the original list by the key returned by the given keySelector function
// applied to each element and returns a map where each group key is associated with a list of corresponding elements.
func GroupBy[K comparable, V any](elements []V, keySelector func(V) K) map[K][]V {
	counts := map[K][]V{}
	for _, element := range elements {
		key := keySelector(element)
		counts[key] = append(counts[key], element)
	}
	return counts
}
