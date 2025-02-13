package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
)

type File struct {
	Name      string
	Content   string
	Size      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Directory struct {
	Name      string
	Files     map[string]*File
	SubDirs   map[string]*Directory
	Parent    *Directory
	CreatedAt time.Time
}

type VFS struct {
	Root       *Directory
	CurrentDir *Directory
}

func newVFS() *VFS {
	root := &Directory{
		Name:      "/",
		Files:     make(map[string]*File),
		SubDirs:   make(map[string]*Directory),
		CreatedAt: time.Now(),
	}
	return &VFS{Root: root, CurrentDir: root}
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

func main() {
	vfs := newVFS()

	commands := map[string]func([]string){
		"cd": func(args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: cd <directory>")
				return
			}
			vfs.cd(args[0])
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
			if len(args) != 1 {
				fmt.Println("Usage: touch <file-name>")
				return
			}
			vfs.touch(args[0])
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

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("&Shell : ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		parts, err := shellwords.Parse(input)
		if err != nil {
			fmt.Println("Error parsing input:", err)
			continue
		}

		if len(parts) == 0 {
			continue
		}

		commandName := parts[0]
		args := parts[1:]

		command, ok := commands[commandName]
		if !ok {
			if commandName == "exit" {
				fmt.Println("Exiting")
				break
			}
			fmt.Println("Unknown command:", commandName)
			continue
		}
		command(args)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
