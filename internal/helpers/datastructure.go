package helpers

func GetKeysFromMap(m *map[string]string) []string {
	// Get Type
	s := make([]string, 0)
	for k, _ := range *m {
		s = append(s, k)
	}

	return s
}
