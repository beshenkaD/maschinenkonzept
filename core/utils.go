package core

func in(s string, a ...string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}

	return false
}
