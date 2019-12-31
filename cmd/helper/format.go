package helper

// MapStringsToSliceString is used to convert map[string]string to a slice of strings, combining
// the key and value via the specified join string.
func MapStringsToSliceString(input map[string]string, join string) []string {
	s := []string{}
	for k, v := range input {
		s = append(s, k+join+v)
	}
	return s
}
