package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
)

func newVFS() *VFS {
	root := &Directory{
		Name:            "/",
		Files:           make(map[string]*File),
		SubDirs:         make(map[string]*Directory),
		CreatedAt:       time.Now(),
		Parent:          "",
		Path:            "/",
		History:         []string{"init"},
		ReadPermission:  []int{-1, 0},
		WritePermission: []int{-1, 0},
	}
	return &VFS{Root: root, CurrentDir: root}
}

func inputs(vfs *VFS, commands CommandMap) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("&Shell : ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if len(input) == 0 {
			continue
		}

		if vfs.CurrentDir != vfs.Root {
			vfs.CurrentDir.History = append(vfs.CurrentDir.History, input)
		} else {
			vfs.Root.History = append(vfs.Root.History, input)
		}

		parts := strings.Split(input, " >> ")
		var commandName string
		var args []string

		if len(parts) > 1 {
			leftSideParts, err := shellwords.Parse(parts[0])
			if err != nil {
				fmt.Println("Error parsing command:", err)
				continue
			}
			rightSide := parts[1]
			commandName = leftSideParts[0]
			args = append(leftSideParts[1:], ">>", rightSide)

		} else {
			parsedParts, err := shellwords.Parse(input)
			if err != nil {
				fmt.Println("Error parsing command:", err)
				continue
			}
			commandName = parsedParts[0]
			args = parsedParts[1:]
		}

		if commandName == "exit" {
			fmt.Println("Exiting")
			err := saveStruct("filedata.gob", vfs)
			if err != nil {
				fmt.Println(err)
			}
			break
		}

		command, ok := commands[commandName]
		if !ok {
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
	vfs, err := loadStruct("filedata.gob")
	if err != nil {
		vfs = newVFS()
		vfs.initAdmin()
	}
	usage := GetUsage()
	commands := GetCommands(vfs, usage)
	inputs(vfs, commands)
}
