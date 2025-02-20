package main

func checkOverlap(arr1, arr2 []int) bool {
	set := make(map[int]struct{})

	for _, num := range arr1 {
		set[num] = struct{}{}
	}

	for _, num := range arr2 {
		if _, exists := set[num]; exists {
			return true
		}
	}
	return false
}

func getIndex(arr1, arr2 []int) (bool, int) {
	set := make(map[int]int)

	for i, num := range arr1 {
		set[num] = i
	}

	for _, num := range arr2 {
		if index, exists := set[num]; exists {
			return true, index
		}
	}
	return false, -1
}
func removeElementByIndex(slice []int, index int) []int {
	sliceLen := len(slice)
	sliceLastIndex := sliceLen - 1
	if index != sliceLastIndex {
		slice[index] = slice[sliceLastIndex]
	}
	return slice[:sliceLastIndex]
}
