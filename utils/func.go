package utils

func Map[T any, Q any](array []T, f func(a T) Q) []Q {
	var result []Q
	for _, a := range array {
		result = append(result, f(a))
	}
	return result
}

func Filter[T any](array []T, f func(a T) bool) []T {
	var result []T
	for _, a := range array {
		if f(a) {
			result = append(result, a)
		}
	}
	return result
}

func All(array []bool) bool {
	for _, item := range array {
		if !item {
			return false
		}
	}
	return true
}

func Any(array []bool) bool {
	for _, item := range array {
		if item {
			return true
		}
	}
	return false
}
