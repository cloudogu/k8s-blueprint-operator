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

// MapWithFunction takes a slice of a type T and a converter function in order to return a slice of type V.
func MapWithFunction[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
