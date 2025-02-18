package main

import (
	"bufio"
	"fmt"
	"os"
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
	vfs.initAdmin()
	commands := GetCommands(vfs)
	inputs(vfs, commands)
}
