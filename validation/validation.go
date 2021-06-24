package validation

func CheckEnum(allowed []string, value string) bool {
	found := false
	for _, a := range allowed {
		if a == value {
			found = true
			break
		}
	}
	return found
}
