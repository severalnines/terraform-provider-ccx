package resources

func parametersEqual(existing, next map[string]string) bool {
	if len(existing) != len(next) {
		return false
	}

	for k, v := range next {
		if e, ok := existing[k]; !ok || e != v {
			return false
		}
	}

	return true
}
