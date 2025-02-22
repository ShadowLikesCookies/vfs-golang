package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

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

func openInEditor(content string) (string, error) {
	tempfile, err := os.CreateTemp("", "temporaryTextFile")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(tempfile.Name(), []byte(content), 0644)
	defer os.Remove(tempfile.Name())

	if err := tempfile.Close(); err != nil {
		return "", err

	}

	cmd := exec.Command("nvim", tempfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Nvim exited with: %w", err)

	}

	editedContent, err := os.ReadFile(tempfile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	defer tempfile.Close()
	defer os.Remove(tempfile.Name())
	return string(editedContent), nil
}
