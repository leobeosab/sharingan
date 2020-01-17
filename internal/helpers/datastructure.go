package helpers

func GetKeysFromMap(m *map[string]string) []string {
	// Get Type
	s := make([]string, 0)
	for k, _ := range *m {
		s = append(s, k)
	}

	return s
}

func RemoveDuplicatesInSlice(slice []string) []string {

	check := make(map[string]int)
	res := make([]string, 0)

	for _, val := range slice {
		check[val] = 0
	}

	for s, _ := range check {
		res = append(res, s)
	}

	return res
}
