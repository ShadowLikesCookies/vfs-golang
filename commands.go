package main

import (
	"fmt"
	"time"
)

func (vfs *VFS) history() {
	for value := range vfs.CurrentDir.History {
		fmt.Println("Value:", vfs.CurrentDir.History[value])
	}
}

func (vfs *VFS) mv(target string, destination string) {
	file, fileExists := vfs.CurrentDir.Files[target]
	dir, dirExists := vfs.CurrentDir.SubDirs[destination]
	if !fileExists {
		fmt.Println("File not found:", target)
		return
	}
	if !dirExists {
		fmt.Println("Destination directory not found:", destination)
		return
	}
	if _, exists := dir.Files[file.Name]; exists {
		fmt.Printf("File %s already exists in %s\n", file.Name, destination)
		return
	}
	dir.Files[file.Name] = file
	delete(vfs.CurrentDir.Files, target)
	fmt.Printf("File %s moved to %s\n", target, destination)
}

func (vfs *VFS) roothistory() {
	for value := range vfs.Root.History {
		fmt.Println("Value:", vfs.Root.History[value])
	}
}
func (vfs *VFS) cd(directory string) {
	if directory == ".." {
		if vfs.CurrentDir.Parent != nil {
			vfs.CurrentDir = vfs.CurrentDir.Parent
		} else {
			fmt.Println("Cannot backtrack, currently at root")
		}
	} else {
		dir, exists := vfs.CurrentDir.SubDirs[directory]
		if !exists {
			fmt.Println("Directory", directory, "does not exist")
			return
		}
		vfs.CurrentDir = dir
	}
}

func (vfs *VFS) ls() {
	for _, file := range vfs.CurrentDir.Files {
		fmt.Println("file:", file.Name)
	}
	for _, dir := range vfs.CurrentDir.SubDirs {
		fmt.Println("dir:", dir.Name)
	}
}

func (vfs *VFS) touch(name string) {
	if _, exists := vfs.CurrentDir.Files[name]; exists {
		fmt.Println("File", name, "already exists")
		return
	}

	file := &File{
		Name:      name,
		Content:   "",
		Size:      0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	vfs.CurrentDir.Files[name] = file
	fmt.Println("File created:", name)
}

func (vfs *VFS) mkdir(name string) {
	if _, exists := vfs.CurrentDir.SubDirs[name]; exists {
		fmt.Println("Directory", name, "already exists")
		return
	}

	dir := &Directory{
		Name:      name,
		Files:     make(map[string]*File),
		SubDirs:   make(map[string]*Directory),
		Parent:    vfs.CurrentDir,
		CreatedAt: time.Now(),
	}
	vfs.CurrentDir.SubDirs[name] = dir
	fmt.Println("Directory created:", name)
}

func (vfs *VFS) pwd() {
	if vfs.CurrentDir.Name == "/" {
		fmt.Println("CWD: /")
	} else {
		fmt.Println("CWD:", vfs.CurrentDir.Name)
	}
}

func (vfs *VFS) cat(name string) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found:", name)
		return
	}
	fmt.Println(file.Content)
}

func (vfs *VFS) fill(amount uint16) {
	for i := uint16(0); i < amount; i++ {
		filename := fmt.Sprintf("file%d", i)
		vfs.touch(filename)
		vfs.mkdir(filename)
	}
}

func (vfs *VFS) echo(name string, content string, appendToFile bool) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		vfs.touch(name)
		file = vfs.CurrentDir.Files[name]
		if file == nil {
			fmt.Println("Error creating file")
			return
		}
	}

	if appendToFile {
		file.Content += content
		fmt.Println("Content appended to file:", name)
	} else {
		file.Content = content
		fmt.Println("Content written to file:", name)
	}

	file.UpdatedAt = time.Now()
	file.Size = len(file.Content)
}

func (vfs *VFS) rm(name string) {
	if _, exists := vfs.CurrentDir.Files[name]; !exists {
		fmt.Println("File not found:", name)
		return
	}
	delete(vfs.CurrentDir.Files, name)
	fmt.Println("File deleted:", name)
}
