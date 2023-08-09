package main

func pushBack[T any](array *[]T, object T) {
	*array = append(*array, object)
}

// func assertEqual[T comparable](expected, actual T) {
// 	if expected != actual {
// 		panic(fmt.Sprintf("Expected %v, got %v", expected, actual))
// 	}
// }

func getUnique[T comparable](array []T) []T {
	seen := make(map[T]bool)
	// input -> map
	for _, item := range array {
		seen[item] = true
	}

	uniqueItems := make([]T, 0, len(seen))
	// map -> output
	for item := range seen {
		uniqueItems = append(uniqueItems, item)
	}

	return uniqueItems
}
