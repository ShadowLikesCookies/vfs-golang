package main

import "fmt"

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

func main2() {
	arr1 := []int{1, 2, 3, 4, 5}
	arr2 := []int{6, 7, 8, 9}

	if checkOverlap(arr1, arr2) {
		fmt.Println("There is an overlap.")
	} else {
		fmt.Println("There is no overlap.")
	}
}
