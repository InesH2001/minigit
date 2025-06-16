package utils

func GetUniqueUnionKeys(maps ...map[string]string) []string {
	keys := make(map[string]struct{})

	for _, m := range maps {
		for k := range m {
			keys[k] = struct{}{}
		}
	}

	var result []string
	for k := range keys {
		result = append(result, k)
	}

	return result
}
