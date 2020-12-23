package utils


// Check if a given string vector has at least one string
func Any(v []string, k string) (bool) {
	for _, str := range(v) {
		if str == k {
			return true
		}
	}

	return false
}