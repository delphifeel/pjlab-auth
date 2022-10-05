package helpers

// efficient remove if order doesn't matter
func SliceRemoveFast[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}