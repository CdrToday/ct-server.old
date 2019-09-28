package main

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func deleteStringFromArray(arr []string, element string) []string {
	for i, e := range arr {
		if e == element {
			arr[i] = arr[len(arr)-1]
			arr = arr[:len(arr)-1]
		}
	}

	return arr
}
