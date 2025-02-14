package main

import "fmt"

func permissionError(error int16) {
	if error == 0 {
		fmt.Println("Error: File dose not have apropriate permissions: read == false")
	}
	if error == 1 {
		fmt.Println("Error: File dose not have apropriate permissions: write == false")
	}
}
