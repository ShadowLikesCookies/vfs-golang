package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-shellwords" // Import the shellwords library
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
		if dir, exists := vfs.CurrentDir.SubDirs[directory]; exists {
			vfs.CurrentDir = dir
		} else {
			fmt.Println("Directory", directory, "does not exist")
		}
	}
}

func (vfs *VFS) touch(name string) {
	file := &File{
		Name:      name,
		Content:   "", // Empty content for touch
		Size:      0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	vfs.CurrentDir.Files[name] = file
	fmt.Println("file created", name)
}

func (vfs *VFS) mkdir(name string) {
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
		fmt.Println("file not found")
		return
	}
	fmt.Println(file.Content)
}

func (vfs *VFS) echo(name string, content string, appendToFile bool) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File not found, creating a new file")
		vfs.touch(name) // Create the file if it doesn't exist
		file = vfs.CurrentDir.Files[name]
		if file == nil {
			fmt.Println("Error creating file")
			return
		}

	}
	if appendToFile {
		file.Content += content
		file.UpdatedAt = time.Now()
		file.Size = len(file.Content)
		fmt.Println("Content appended to file:", name)
	} else {
		file.Content = content
		file.UpdatedAt = time.Now()
		file.Size = len(content)
		fmt.Println("Content written to file:", name)
	}
}

func (vfs *VFS) rm(name string) {
	if _, exists := vfs.CurrentDir.Files[name]; exists {
		delete(vfs.CurrentDir.Files, name)
		fmt.Println("file has been deleted", name)
	} else {
		fmt.Println("File not Found")
	}
}

func main() {
	vfs := newVFS()

	commands := map[string]func([]string){
		"cd": func(args []string) {
			if len(args) > 0 {
				vfs.cd(args[0])
			} else {
				fmt.Println("Usage: cd <directory>")
			}
		},
		"pwd": func(args []string) {
			vfs.pwd() // Pwd doesn't take args, so don't check length
		},
		"rm": func(args []string) {
			if len(args) > 0 {
				vfs.rm(args[0])
			} else {
				fmt.Println("Usage: rm <file-name>")
			}
		},
		"mkdir": func(args []string) {
			if len(args) > 0 {
				vfs.mkdir(args[0])
			} else {
				fmt.Println("Usage: mkdir <dir-name>")
			}
		},
		"touch": func(args []string) {
			if len(args) > 0 {
				vfs.touch(args[0])
			} else {
				fmt.Println("Usage: touch <file-name>")
			}

		},
		"echo": func(args []string) {
			if len(args) > 1 {
				filename := args[0]
				content := strings.Join(args[1:], " ")
				vfs.echo(filename, content, false) //Over write is false.
			} else {
				fmt.Println("Usage: echo <file-name> <content>")
			}
		},
		"cat": func(args []string) {
			if len(args) > 0 {
				vfs.cat(args[0])
			} else {
				fmt.Println("Usage: cat <file-name>")
			}
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("&Shell : ")
		scanner.Scan()
		input := scanner.Text()

		// Use shellwords.Parse to handle quotes
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

		if command, ok := commands[commandName]; ok {
			command(args)
		} else if commandName == "exit" {
			fmt.Println("Exiting")
			break
		} else {
			fmt.Println("Unknown command", commandName)
		}

	}
}
