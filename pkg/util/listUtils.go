package util

func GetDuplicates[T comparable](list []T) []T {
	elementCount := make(map[T]int)

	// countByValue
	for _, value := range list {
		elementCount[value] += 1
	}

	// get list of values with count != 1
	var duplicates []T
	for value, count := range elementCount {
		if count != 1 {
			duplicates = append(duplicates, value)
		}
	}
	return duplicates
}

// Map returns a slice where each element is the result of invoking fn on each corresponding element of the given slice.
//
// Example:
//
//	fruits := []string{"apple", "banana", "raspberry"}
//	loudFruits := Map(fruits, strings.ToUpper)
//	fmt.Println(loudFruits)
//
// This should print: [APPLE BANANA RASPBERRY]
func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
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
