package collectons

func IndexOfStringSlice(s []string, target string) int {
	for i := range s {
		if s[i] == target {
			return i
		}
	}
	return -1
}
