package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetUsage() UsageMap {
	return UsageMap{
		"cd": func() {
			fmt.Println("Usage: cd <directory>")
		},
		"mv": func() {
			fmt.Println("Usage: mv <target> <destination>")
		},
		"history": func() {
			fmt.Println("Usage: history")
		},
		"roothistory": func() {
			fmt.Println("Usage: roothistory")
		},
		"pwd": func() {
			fmt.Println("Usage: pwd")
		},
		"rm": func() {
			fmt.Println("Usage: rm <file-name>")
		},
		"ls": func() {
			fmt.Println("Usage: ls")
		},
		"fill": func() {
			fmt.Println("Usage: fill <amount>")
		},
		"mkdir": func() {
			fmt.Println("Usage: mkdir <dir-name>")
		},
		"touch": func() {
			fmt.Println("Usage: touch <file-name>")
		},
		"echo": func() {
			fmt.Println("Usage: echo <file-name> <content>")
		},
		"cat": func() {
			fmt.Println("Usage: cat <file-name> [>> <destination-file>]")
		},
		"remPerms": func() {
			fmt.Println("Usage: remPerms <file-name> <permission> <id>")
		},
		"whoami": func() {
			fmt.Println("Usage: whoami")
		},
		"addPerms": func() {
			fmt.Println("Usage: addPerms <file-name> <permission> <id>")
		},
		"nvim": func() {
			fmt.Println("Usage: nvim <file-name> | nvim . ")
		},
		"clear": func() {
			fmt.Println("Usage: clear")
		},
		"call": func() {
			fmt.Println("Usage: call <name>")
		},
	}
}

func GetCommands(vfs *VFS, usage UsageMap) CommandMap {
	return CommandMap{
		"cd": func(args []string) {
			if len(args) != 1 {
				usage["cd"]()
				return
			}
			vfs.cd(args[0])
		},
		"mv": func(args []string) {
			if len(args) != 2 {
				usage["mv"]()
				return
			}
			vfs.mv(args[0], args[1])
			fmt.Println("Moved", args[0], "to", args[1])
		},
		"history": func(args []string) {
			if len(args) != 0 {
				usage["history"]()
				return
			}

			vfs.history()
			fmt.Println("Displayed history")
		},
		"roothistory": func(args []string) {
			if len(args) != 0 {
				usage["roothistory"]()
				return
			}
			vfs.roothistory()
			fmt.Println("Displayed root history")
		},
		"pwd": func(args []string) {
			if len(args) != 0 {
				usage["pwd"]()
				return
			}
			vfs.pwd()
		},
		"rm": func(args []string) {
			if len(args) != 1 {
				usage["rm"]()
				return
			}
			vfs.rm(args[0])
			fmt.Println("Removed file", args[0])
		},
		"ls": func(args []string) {
			if len(args) != 0 {
				usage["ls"]()
				return
			}
			vfs.ls()
			fmt.Println("Listed directory contents")
		},
		"fill": func(args []string) {
			if len(args) != 1 {
				usage["fill"]()
				return
			}
			amount, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				return
			}
			vfs.fill(uint16(amount))
			fmt.Println("Filled directory with", amount, "files and directories")
		},
		"mkdir": func(args []string) {
			if len(args) != 1 {
				usage["mkdir"]()
				return
			}
			vfs.mkdir(args[0])
			fmt.Println("Created directory", args[0])
		},
		"touch": func(args []string) {
			if len(args) != 1 {
				usage["touch"]()
				return
			}
			vfs.touch(args[0])
		},
		"echo": func(args []string) {
			if len(args) < 2 {
				usage["echo"]()
				return
			}
			vfs.echo(args[0], strings.Join(args[1:], ""), false)
			fmt.Println("Written to file", args[0])
		},
		"cat": func(args []string) {
			if len(args) == 1 {
				contentPtr := vfs.cat(args[0])
				fmt.Println("Content: ", *contentPtr)
			} else if len(args) == 3 && args[1] == ">>" {
				sourceFileName := args[0]
				destFileName := args[2]

				contentPtr := vfs.cat(sourceFileName)

				destFile, exists := vfs.CurrentDir.Files[destFileName]
				if !exists {
					fmt.Println("Destination file not found:", destFileName)
					return
				}

				err := vfs.pipe(contentPtr, destFile)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Piped content to ", destFileName)

			} else {
				usage["cat"]()
				return
			}
		},

		"remPerms": func(args []string) {
			if len(args) != 3 {
				usage["remPerms"]()
				return
			}
			intstringconverted, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			vfs.remPerms(args[0], args[1], intstringconverted)
			fmt.Println("Removed permission", args[1], "from", args[0], "for ID", args[2])
		},
		"whoami": func(args []string) {
			if len(args) != 0 {
				usage["whoami"]()
				return
			} else {
				fmt.Println("Current User: ", *vfs.whoami())
			}
		},
		"addPerms": func(args []string) {
			if len(args) != 3 {
				usage["addPerms"]()
				return
			}
			stringintconverted, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			vfs.addPerms(args[0], args[1], stringintconverted)
			fmt.Println("Added permission", args[1], "to", args[0], "for ID", args[2])
		},
		"nvim": func(args []string) {
			if len(args) != 1 {
				usage["nvim"]()
				return
			}
			vfs.nvim(args[0])
		},
		"clear": func(args []string) {
			if len(args) != 0 {
				usage["clear"]()
				return
			}
			vfs.clear()
		},
		"call": func(args []string) {
			if len(args) != 1 {
				usage["call"]()
				return
			}
			vfs.call(args[0])
		},
		"time": func(args []string) {
			vfs.date()
		},
	}
}

func (vfs *VFS) initAdmin() {
	user := &User{
		Name:       "admin",
		GroupPerms: []int{0, -1},
	}
	vfs.CurrentUser = user
}
func (vfs *VFS) clear() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
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

	if !checkOverlap(vfs.CurrentDir.SubDirs[destination].WritePermission, vfs.CurrentUser.GroupPerms) {
		fmt.Println("You do not have write permissions in the destination directory.")
		return
	}
	if !checkOverlap(file.WritePermission, vfs.CurrentUser.GroupPerms) {
		fmt.Println("You do not have write permissions to move this file.")
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
		if vfs.CurrentDir.Parent == "" {
			fmt.Println("Cannot backtrack, currently at root")
			return
		} else {
			parentDir := vfs.findDirectoryByPath(vfs.CurrentDir.Parent)
			if parentDir != nil {
				vfs.CurrentDir = parentDir
			} else {
				fmt.Println("Error: Parent directory not found!")
				return
			}
		}
	} else {
		dir, exists := vfs.CurrentDir.SubDirs[directory]
		if !exists {
			fmt.Println("Directory", directory, "does not exist")
			return
		}
		if !checkOverlap(dir.ReadPermission, vfs.CurrentUser.GroupPerms) {
			fmt.Println("You do not have read permissions to access this directory.")
			return
		}
		vfs.CurrentDir = dir
	}
}

func (vfs *VFS) call(name string) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File ", name, "Dose not exist")
		return
	}
	if !file.Executable {
		fmt.Println("File", file.Name, "Dose not have Executable permissions")
		return
	}
	parts := strings.Split(file.Name, ".")
	if parts[len(parts)-1] == "vsh" {
		vfs.executeArray(vfs.getCommandArray(name))
	}
}

func (vfs *VFS) ls() (filearray []string, dirarray []string) {
	if !checkOverlap(vfs.CurrentDir.ReadPermission, vfs.CurrentUser.GroupPerms) {
		fmt.Println("You do not have read permissions to list this directory.")
		return filearray, dirarray
	}

	for _, file := range vfs.CurrentDir.Files {
		filearray = append(filearray, file.Name)
		fmt.Println("file:", file.Name)
	}
	for _, dir := range vfs.CurrentDir.SubDirs {
		dirarray = append(dirarray, dir.Name)
		fmt.Println("dir:", dir.Name)
	}
	return filearray, dirarray
}

func (vfs *VFS) touch(name string) {
	if !checkOverlap(vfs.CurrentDir.WritePermission, vfs.CurrentUser.GroupPerms) {
		fmt.Println("You do not have write permissions to create files in this directory.")
		return
	}
	if _, exists := vfs.CurrentDir.Files[name]; exists {
		fmt.Println("File", name, "already exists")
		return
	}
	pattern := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*\.[a-z0-9]+$`)
	isvalid := pattern.MatchString(name)
	if !isvalid {
		fmt.Println("Name dose not match expected format <1-9,a-z.1-9.a-z")
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
		Executable:       false,
	}
	vfs.CurrentDir.Files[name] = file
	fmt.Println("Created file", name)
}

func (vfs *VFS) nvim(name string) {

	if name == "." && checkOverlap(vfs.CurrentUser.GroupPerms, vfs.CurrentDir.ReadPermission) {
		arr1, arr2 := vfs.ls()
		arr1 = append([]string{"Files: "}, arr1...)
		arr1 = append(arr1, "\n")
		arr2 = append([]string{"Directories: "}, arr2...)
		arr2 = append(arr2, "\n")
		combined := append(arr1, arr2...)
		_, err := openInEditor(strings.Join(combined, "\n"), false)
		if err != nil {
			fmt.Println("Error has occured whilst open nvim %w", err)
		}
		return
	}

	if _, exists := vfs.CurrentDir.Files[name]; !exists {
		vfs.touch(name)
	}
	if !checkOverlap(vfs.CurrentUser.GroupPerms, vfs.CurrentDir.Files[name].ReadPermission) {
		fmt.Println("You do not have the apropriate Read permissions")
		return
	}
	editedText, err := openInEditor(vfs.CurrentDir.Files[name].Content, true)
	if err != nil {
		fmt.Println("Error has occured whilst open nvim %w", err)
		return
	}

	if !checkOverlap(vfs.CurrentUser.GroupPerms, vfs.CurrentDir.Files[name].WritePermission) {
		fmt.Println("You do not have the apropriate Write permissions")
		return
	}
	vfs.CurrentDir.Files[name].Content = *editedText
}

func (vfs *VFS) mkdir(name string) {
	if !checkOverlap(vfs.CurrentDir.WritePermission, vfs.CurrentUser.GroupPerms) {
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
		Parent:           vfs.CurrentDir.Path,
		CreatedAt:        time.Now(),
		Path:             vfs.CurrentDir.Path + "/" + name,
		ReadPermission:   []int{1, -1},
		WritePermission:  []int{1, -1},
		ModifyPermission: []int{1, -1},
	}

	vfs.CurrentDir.SubDirs[name] = dir
	fmt.Println("Directory created:", name)
}

func (vfs *VFS) addPerms(name string, permission string, id int) {
	permission = strings.ToLower(permission)
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found:", name)
		return
	} else {
		if checkOverlap(vfs.CurrentUser.GroupPerms, vfs.CurrentDir.Files[name].ModifyPermission) {
			if permission == "write" {
				vfs.CurrentDir.Files[name].WritePermission = append(vfs.CurrentDir.Files[name].WritePermission, id)
			} else if permission == "read" {
				vfs.CurrentDir.Files[name].ReadPermission = append(vfs.CurrentDir.Files[name].ReadPermission, id)
			} else if permission == "modify" {
				vfs.CurrentDir.Files[name].ModifyPermission = append(file.ModifyPermission, id)
			} else if permission == "executable" {
				if id == 0 {
					vfs.CurrentDir.Files[name].Executable = false
				} else if id == 1 {
					vfs.CurrentDir.Files[name].Executable = true
				} else {
					fmt.Println("For executable permission, value must be between 0-1")
					return
				}
			} else {
				fmt.Println("Permission dose not exist")
			}
		}
	}
}

func (vfs *VFS) remPerms(name string, permission string, id int) {
	permission = strings.ToLower(permission)
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found:", name)
		return
	} else {
		if checkOverlap(vfs.CurrentUser.GroupPerms, vfs.CurrentDir.Files[name].ModifyPermission) {
			if permission == "write" {
				exists, index := getIndex(file.WritePermission, []int{id})
				if exists {
					temp := removeElementByIndex(file.WritePermission, index)
					file.WritePermission = temp
					temp = nil
				} else {
					fmt.Println("Permission ID dose not exist in writePermissions[]")
				}
			} else if permission == "read" {
				exists, index := getIndex(file.ReadPermission, []int{id})
				if exists {
					temp := removeElementByIndex(file.ReadPermission, index)
					file.ReadPermission = temp
					temp = nil
				} else {
					fmt.Println("Permission ID dose not exist in ReadPermissions[]")
				}
			} else if permission == "modify" {
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
	if checkOverlap(vfs.CurrentDir.Files[name].ReadPermission, vfs.CurrentUser.GroupPerms) {
		return &file.Content
	} else {
		fmt.Println("You do not share any permission ID's with this file. READ==FALSE")
	}
	return nil
}

func (vfs *VFS) fill(amount uint16) {
	if !checkOverlap(vfs.CurrentDir.WritePermission, vfs.CurrentUser.GroupPerms) {
		fmt.Println("You do not have write permissions to fill this directory.")
		return
	}
	for i := uint16(0); i < amount; i++ {
		filename := fmt.Sprintf("file%d", i)
		vfs.touch(filename + ".txt")
		vfs.mkdir(filename + ".txt")
	}
}

func (vfs *VFS) echo(name string, content string, appendToFile bool) {
	file, exists := vfs.CurrentDir.Files[name]
	if exists {
		if !checkOverlap(file.WritePermission, vfs.CurrentUser.GroupPerms) {
			fmt.Println("You do not share any group permissions to WRITE to this file.")
			return
		}
	} else {
		if !checkOverlap(vfs.CurrentDir.WritePermission, vfs.CurrentUser.GroupPerms) {
			fmt.Println("You do not have write permissions to create files in this directory.")
			return
		}
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
}

func (vfs *VFS) whoami() *string {
	return &vfs.CurrentUser.Name
}

func (vfs *VFS) rm(name string) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found:", name)
		return
	}
	if !checkOverlap(file.WritePermission, vfs.CurrentUser.GroupPerms) {
		fmt.Println("You do not have write permissions to delete this file.")
		return
	}
	delete(vfs.CurrentDir.Files, name)
	fmt.Println("File deleted:", name)
}

func (vfs *VFS) date() {
	fmt.Println("Current Time: ", time.Now())
}
