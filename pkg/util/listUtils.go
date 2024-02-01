package util

func GetDuplicates(list []string) []string {
	elementCount := make(map[string]int)

	// countByName
	for _, name := range list {
		elementCount[name] += 1
	}

	// get list of names with count != 1
	var duplicates []string
	for name, count := range elementCount {
		if count != 1 {
			duplicates = append(duplicates, name)
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
