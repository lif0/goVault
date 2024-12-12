package btypes

// To return an array with one element.
func ToArray[T any](t T) []T {
	arr := make([]T, 1)
	arr[0] = t
	return arr
}
