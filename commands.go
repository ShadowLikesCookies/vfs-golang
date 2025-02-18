package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CommandMap map[string]func([]string)

func GetCommands(vfs *VFS) CommandMap {
	return CommandMap{
		"cd": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: cd <directory>")
				return
			}
			vfs.cd(args[0])
		},
		"mv": func(args []string) {
			if len(args) < 2 {
				fmt.Println("Usage: mv <target> <destination>")
			}
			vfs.mv(args[0], args[1])
		},
		"history": func(args []string) {
			vfs.history()
		},
		"roothistory": func(args []string) {
			vfs.roothistory()
		},
		"pwd": func(args []string) {
			vfs.pwd()
		},
		"rm": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: rm <file-name>")
				return
			}
			vfs.rm(args[0])
		},
		"ls": func(args []string) {
			vfs.ls()
		},
		"fill": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: fill <amount>")
				return
			}
			amount, err := strconv.ParseUint(args[0], 10, 16)
			if err != nil {
				fmt.Println("Error converting string to uint16:", err)
				return
			}
			vfs.fill(uint16(amount))
		},
		"mkdir": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: mkdir <dir-name>")
				return
			}
			vfs.mkdir(args[0])
		},
		"touch": func(args []string) {
			if len(args) != 3 {
				fmt.Println("Usage: touch <file-name> <read-perm> <write-perm>")
				return
			}

			fileName := args[0]
			readPermStr := args[1]
			writePermStr := args[2]
			readPerm, err := strconv.ParseBool(readPermStr)
			if err != nil {
				fmt.Println("Invalid read permission value.  Must be true or false:", err)
				return
			}

			writePerm, err := strconv.ParseBool(writePermStr)
			if err != nil {
				fmt.Println("Invalid write permission value. Must be true or false:", err)
				return
			}

			vfs.touch(fileName, []bool{readPerm, writePerm})
		},
		"echo": func(args []string) {
			if len(args) < 2 {
				fmt.Println("Usage: echo <file-name> <content>")
				return
			}
			filename := args[0]
			content := strings.Join(args[1:], " ")
			vfs.echo(filename, content, false)
		},
		"cat": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: cat <file-name>")
				return
			}
			vfs.cat(args[0])
		},
	}
}

func (vfs *VFS) initAdmin() {
	user := &User{
		name:       "admin",
		groupPerms: []int16{0},
	}
	vfs.CurrentUser = user
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
func (vfs *VFS) history() {
	for value := range vfs.CurrentDir.History {
		fmt.Println("Value:", vfs.CurrentDir.History[value])
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

func (vfs *VFS) touch(name string, Permissions []bool) {
	if _, exists := vfs.CurrentDir.Files[name]; exists {
		fmt.Println("File", name, "already exists")
		return
	}

	file := &File{
		Name:        name,
		Content:     "",
		Size:        0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Permissions: []bool{Permissions[0], Permissions[1], false},
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
	if !vfs.CurrentDir.Files[name].Permissions[0] == false {
		fmt.Println(file.Content)
	} else {
		permissionError(0)
	}
}

func (vfs *VFS) fill(amount uint16) {
	for i := uint16(0); i < amount; i++ {
		filename := fmt.Sprintf("file%d", i)
		vfs.touch(filename, []bool{true})
		vfs.mkdir(filename)
	}
}

func (vfs *VFS) echo(name string, content string, appendToFile bool) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		vfs.touch(name, []bool{true, true})
		file = vfs.CurrentDir.Files[name]
		if file == nil {
			fmt.Println("Error creating file")
			return
		}
	}
	if vfs.CurrentDir.Files[name].Permissions[1] != false {
		if appendToFile {
			file.Content += content
			fmt.Println("Content appended to file:", name)
		} else {
			file.Content = content
			fmt.Println("Content written to file:", name)
		}

		file.UpdatedAt = time.Now()
		file.Size = len(file.Content)
	} else {
		permissionError(1)
	}

}

func (vfs *VFS) rm(name string) {
	if _, exists := vfs.CurrentDir.Files[name]; !exists {
		fmt.Println("File not found:", name)
		return
	}
	delete(vfs.CurrentDir.Files, name)
	fmt.Println("File deleted:", name)
}
