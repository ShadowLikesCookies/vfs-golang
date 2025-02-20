package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetCommands(vfs *VFS) CommandMap {
	return CommandMap{
		"cd": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: cd <directory>")
				return
			}
			vfs.cd(args[0])
			fmt.Println("Changed directory to", args[0])
		},
		"mv": func(args []string) {
			if len(args) < 2 {
				fmt.Println("Usage: mv <target> <destination>")
				return
			}
			vfs.mv(args[0], args[1])
			fmt.Println("Moved", args[0], "to", args[1])
		},
		"history": func(args []string) {
			vfs.history()
			fmt.Println("Displayed history")
		},
		"roothistory": func(args []string) {
			vfs.roothistory()
			fmt.Println("Displayed root history")
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
			fmt.Println("Removed file", args[0])
		},
		"ls": func(args []string) {
			vfs.ls()
			fmt.Println("Listed directory contents")
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
			fmt.Println("Filled directory with", amount, "files and directories")
		},
		"mkdir": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: mkdir <dir-name>")
				return
			}
			vfs.mkdir(args[0])
			fmt.Println("Created directory", args[0])
		},
		"touch": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: touch <file-name>")
				return
			}
			fileName := args[0]
			vfs.touch(fileName)
			fmt.Println("Created file", fileName)
		},
		"echo": func(args []string) {
			if len(args) < 2 {
				fmt.Println("Usage: echo <file-name> <content>")
				return
			}
			filename := args[0]
			content := strings.Join(args[1:], " ")
			vfs.echo(filename, content, false)
			fmt.Println("Written to file", filename)
		},
		"cat": func(args []string) {
			if len(args) < 1 {
				fmt.Println("Usage: cat <file-name> [>> <destination-file>]")
				return
			}

			sourceFileName := args[0]
			if len(args) > 1 && args[1] == ">>" {
				if len(args) < 3 {
					fmt.Println("Usage: cat <file-name> >> <destination-file>")
					return
				}
				destFileName := args[2]
				contentPtr := vfs.cat(sourceFileName)
				if contentPtr != nil {
					vfs.echo(destFileName, *contentPtr, true)
				}
			} else {
				vfs.cat(sourceFileName)
				fmt.Println("Content: ", vfs.CurrentDir.Files[sourceFileName].Content)
			}
		},
		"remPerms": func(args []string) {
			if len(args) != 3 {
				fmt.Println("Usage: remPerms <file-name> <permission> <id>")
				return
			}
			tempInt64, err := strconv.ParseInt(args[2], 0, 0)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			vfs.remPerms(args[0], args[1], int(tempInt64))
			fmt.Println("Removed permission", args[1], "from", args[0], "for ID", args[2])
		},
		"addPerms": func(args []string) {
			if len(args) != 3 {
				fmt.Println("Usage: addPerms <file-name> <permission> <id>")
				return
			}
			tempInt64, err := strconv.ParseInt(args[2], 0, 0)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			vfs.addPerms(args[0], args[1], int(tempInt64))
			fmt.Println("Added permission", args[1], "to", args[0], "for ID", args[2])
		},
	}
}

func (vfs *VFS) initAdmin() {
	user := &User{
		name:       "admin",
		groupPerms: []int{0, -1},
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

	if !checkOverlap(file.WritePermission, vfs.CurrentUser.groupPerms) {
		fmt.Println("You do not have write permissions to move this file.")
		return
	}

	if !checkOverlap(vfs.CurrentDir.SubDirs[destination].WritePermission, vfs.CurrentUser.groupPerms) {
		fmt.Println("You do not have write permissions in the destination directory.")
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
		if !checkOverlap(dir.ReadPermission, vfs.CurrentUser.groupPerms) {
			fmt.Println("You do not have read permissions to access this directory.")
			return
		}
		vfs.CurrentDir = dir
	}
}

func (vfs *VFS) ls() {
	if !checkOverlap(vfs.CurrentDir.ReadPermission, vfs.CurrentUser.groupPerms) {
		fmt.Println("You do not have read permissions to list this directory.")
		return
	}

	for _, file := range vfs.CurrentDir.Files {
		fmt.Println("file:", file.Name)
	}
	for _, dir := range vfs.CurrentDir.SubDirs {
		fmt.Println("dir:", dir.Name)
	}
}

func (vfs *VFS) touch(name string) {
	if !checkOverlap(vfs.CurrentDir.WritePermission, vfs.CurrentUser.groupPerms) {
		fmt.Println("You do not have write permissions to create files in this directory.")
		return
	}
	if _, exists := vfs.CurrentDir.Files[name]; exists {
		fmt.Println("File", name, "already exists")
		return
	}

	file := &File{
		Name:             name,
		Content:          "",
		Size:             0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ReadPermission:   []int{1, -1},
		WritePermission:  []int{1, -1},
		ModifyPermission: []int{1, -1},
	}
	vfs.CurrentDir.Files[name] = file
}

func (vfs *VFS) mkdir(name string) {
	if !checkOverlap(vfs.CurrentDir.WritePermission, vfs.CurrentUser.groupPerms) {
		fmt.Println("You do not have write permissions to create directories in this directory.")
		return
	}
	if _, exists := vfs.CurrentDir.SubDirs[name]; exists {
		fmt.Println("Directory", name, "already exists")
		return
	}

	dir := &Directory{
		Name:             name,
		Files:            make(map[string]*File),
		SubDirs:          make(map[string]*Directory),
		Parent:           vfs.CurrentDir,
		CreatedAt:        time.Now(),
		ReadPermission:   []int{1, -1},
		WritePermission:  []int{1, -1},
		ModifyPermission: []int{1, -1},
	}

	vfs.CurrentDir.SubDirs[name] = dir
	fmt.Println("Directory created:", name)
}

func (vfs *VFS) addPerms(name string, permission string, id int) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found:", name)
		return
	} else {
		if checkOverlap(vfs.CurrentUser.groupPerms, vfs.CurrentDir.Files[name].ModifyPermission) {
			if strings.ToLower(permission) == "write" {
				vfs.CurrentDir.Files[name].WritePermission = append(vfs.CurrentDir.Files[name].WritePermission, id)
			} else if strings.ToLower(permission) == "read" {
				vfs.CurrentDir.Files[name].ReadPermission = append(vfs.CurrentDir.Files[name].ReadPermission, id)
			} else if strings.ToLower(permission) == "modify" {
				vfs.CurrentDir.Files[name].ModifyPermission = append(file.ModifyPermission, id)
			} else {
				fmt.Println("Permission dose not exist")
			}
		}
	}
}

func (vfs *VFS) remPerms(name string, permission string, id int) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found:", name)
		return
	} else {
		if checkOverlap(vfs.CurrentUser.groupPerms, vfs.CurrentDir.Files[name].ModifyPermission) {
			if strings.ToLower(permission) == "write" {
				exists, index := getIndex(file.WritePermission, []int{id})
				if exists {
					temp := removeElementByIndex(file.WritePermission, index)
					file.WritePermission = temp
					temp = nil
				} else {
					fmt.Println("Permission ID dose not exist in writePermissions[]")
				}
			} else if strings.ToLower(permission) == "read" {
				exists, index := getIndex(file.ReadPermission, []int{id})
				if exists {
					temp := removeElementByIndex(file.ReadPermission, index)
					file.ReadPermission = temp
					temp = nil
				} else {
					fmt.Println("Permission ID dose not exist in ReadPermissions[]")
				}
			} else if strings.ToLower(permission) == "modify" {
				exists, index := getIndex(file.ModifyPermission, []int{id})
				if exists {
					temp := removeElementByIndex(file.ModifyPermission, index)
					file.ModifyPermission = temp
					temp = nil
				} else {
					fmt.Println("Permission ID dose not exist in ModifyPermissions[]")
				}
			}
		}
	}
}

func (vfs *VFS) pwd() {
	fmt.Println("CWD:", vfs.CurrentDir.Name)
}

func (vfs *VFS) cat(name string) *string {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found:", name)
		return nil
	}
	if checkOverlap(vfs.CurrentDir.Files[name].ReadPermission, vfs.CurrentUser.groupPerms) {
		return &file.Content
	} else {
		fmt.Println("You do not share any permission ID's with this file. READ==FALSE")
	}
	return nil
}

func (vfs *VFS) fill(amount uint16) {
	if !checkOverlap(vfs.CurrentDir.WritePermission, vfs.CurrentUser.groupPerms) {
		fmt.Println("You do not have write permissions to fill this directory.")
		return
	}
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

	if checkOverlap(vfs.CurrentDir.Files[name].WritePermission, vfs.CurrentUser.groupPerms) {
		if appendToFile {
			file.Content += content
			fmt.Println("Content appended to file:", name)
		} else {
			file.Content = content
			fmt.Println("Content written to file:", name)
		}
	} else {
		fmt.Println("You do not share any group permissions. WRITE==FALSE")
	}
}

func (vfs *VFS) rm(name string) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found:", name)
		return
	}
	if !checkOverlap(file.WritePermission, vfs.CurrentUser.groupPerms) {
		fmt.Println("You do not have write permissions to delete this file.")
		return
	}
	delete(vfs.CurrentDir.Files, name)
	fmt.Println("File deleted:", name)
}
