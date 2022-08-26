package common

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if item == v {
			return true
		}
	}

	return false
}

// first match of item in []
func Index(slice []string, item string) int {
	for k, v := range slice {
		if item == v {
			return k
		}
	}

	return -1
}

func RemoveFrom(slice []string, item string) []string {
	var newList []string
	for _, v := range slice {
		if item != v {
			newList = append(newList, v)
		}
	}

	return newList
}
