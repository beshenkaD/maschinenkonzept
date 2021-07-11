package core

func in(s string, a ...string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}

	return false
}

func boolToRus(b bool) string {
	if b {
		return "Да"
	} else {
		return "Нет"
	}
}
