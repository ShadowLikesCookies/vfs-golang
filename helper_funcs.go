package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func loadFromFile(filename string) (*VFS, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var vfs VFS
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&vfs); err != nil {
		return nil, err
	}

	os.Remove("filedata.gob")
	return &vfs, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func (vfs *VFS) findDirectoryByPath(path string) *Directory {
	if path == "/" || path == "" {
		return vfs.Root
	}

	parts := strings.Split(path, "/")
	current := vfs.Root

	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue // Skip empty parts
		}
		nextDir, exists := current.SubDirs[parts[i]]
		if !exists {
			return nil // Directory not found
		}
		current = nextDir
	}
	return current
}
func saveStruct(filename string, data *VFS) error {

	var RealData HelperVFS
	RealData.Root = data.Root
	RealData.CurrentDir = data.CurrentDir
	RealData.CurrentUser = data.CurrentUser

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(RealData); err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}
	return nil
}

func loadStruct(filename string) (*HelperVFS, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	var TempVFS HelperVFS
	if err := decoder.Decode(&TempVFS); err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}
	fmt.Println(&TempVFS)
	return &TempVFS, nil
}

func (vfs *VFS) checkFileWrite(name string) bool {
	return checkOverlap(vfs.CurrentDir.Files[name].WritePermission, vfs.CurrentUser.GroupPerms)
}
func (vfs *VFS) checkFileRead(name string) bool {
	return checkOverlap(vfs.CurrentDir.Files[name].ReadPermission, vfs.CurrentUser.GroupPerms)
}
func (vfs *VFS) checkFileModify(name string) bool {
	return checkOverlap(vfs.CurrentDir.Files[name].ModifyPermission, vfs.CurrentUser.GroupPerms)
}
func (vfs *VFS) checkFileExecutable(name string) bool {
	if vfs.CurrentDir.Files[name].Executable {
		return true
	} else {
		return false
	}
}
func (vfs *VFS) checkDirWrite(name string) bool {
	return checkOverlap(vfs.CurrentDir.SubDirs[name].WritePermission, vfs.CurrentUser.GroupPerms)
}
func (vfs *VFS) checkDirRead(name string) bool {
	return checkOverlap(vfs.CurrentDir.SubDirs[name].ReadPermission, vfs.CurrentUser.GroupPerms)
}
func (vfs *VFS) checkDirModify(name string) bool {
	return checkOverlap(vfs.CurrentDir.SubDirs[name].ModifyPermission, vfs.CurrentUser.GroupPerms)
}

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

func openInEditor(content string, ret bool) (*string, error) {
	tempfile, err := os.CreateTemp("", "temporaryTextFile")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempfile.Name())

	_, err = tempfile.Write([]byte(content))
	if err != nil {
		tempfile.Close()
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}

	if err := tempfile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	cmd := exec.Command("nvim", tempfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("nvim exited with: %w", err)
	}

	editedContent, err := os.ReadFile(tempfile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if ret {
		contentStr := string(editedContent)
		return &contentStr, nil
	} else {
		return nil, nil
	}

}

func (vfs *VFS) getCommandArray(name string) []string {
	fmt.Println(strings.Split(vfs.CurrentDir.Files[name].Content, ";"))
	return strings.Split(vfs.CurrentDir.Files[name].Content, ";")
}
func (vfs *VFS) executeArray(array []string) {
	usage := GetUsage()
	commands := GetCommands(vfs, usage)
	for i := 0; i < len(array); i++ {
		execute(vfs, commands, array[i])
	}
}

func (vfs *VFS) pipe(source *string, file *File) error {
	if source == nil {
		return fmt.Errorf("source data must not be a nil pointer")
	}

	if file == nil {
		return fmt.Errorf("destination *File pointer cannot be nil")
	}

	if !checkOverlap(vfs.CurrentUser.GroupPerms, file.WritePermission) {
		return fmt.Errorf("you do not have the appropriate write permissions for file: %s", file.Name)
	}

	file.Content = *source
	return nil
}
