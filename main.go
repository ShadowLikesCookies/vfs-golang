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

func newVFS() *VFS {
	root := &Directory{
		Name:      "/",
		Files:     make(map[string]*File),
		SubDirs:   make(map[string]*Directory),
		CreatedAt: time.Now(),
		History:   []string{"init"},
	}
	return &VFS{Root: root, CurrentDir: root}
}

func inputs(vfs *VFS, commands map[string]func([]string)) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("&Shell : ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if vfs.CurrentDir != vfs.Root {
			vfs.CurrentDir.History = append(vfs.CurrentDir.History, input)
		} else {
			vfs.Root.History = append(vfs.Root.History, input)
		}

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
	inputs(vfs, commands)
}
