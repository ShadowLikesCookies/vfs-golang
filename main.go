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

func execute(vfs *VFS, commands CommandMap, icommand string) {

	parts := strings.Split(icommand, " >> ")
	var commandName string
	var args []string

	if len(parts) > 1 {
		leftSideParts, err := shellwords.Parse(parts[0])
		if err != nil {
			fmt.Println("Error parsing command:", err)
		}
		rightSide := parts[1]
		commandName = leftSideParts[0]
		args = append(leftSideParts[1:], ">>", rightSide)

	} else {
		parsedParts, err := shellwords.Parse(icommand)
		if err != nil {
			fmt.Println("Error parsing command:", err)
			return
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
		os.Exit(0)
	}

	command, ok := commands[commandName]
	if !ok {
		fmt.Println("Unknown command:", commandName)
		return
	}

	command(args)
}

func inputs(vfs *VFS, commands CommandMap) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("&Shell" + vfs.CurrentDir.Path + ": ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if len(input) == 0 {
			continue
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}
		execute(vfs, commands, input)
	}
}

func main() {
	var vfs *VFS
	TempVFS, err := loadStruct("filedata.gob")
	if err != nil {
		fmt.Println(err)
		vfs = newVFS()
		vfs.initAdmin()
	} else {
		vfs = &VFS{
			Root:        TempVFS.Root,
			CurrentDir:  TempVFS.CurrentDir,
			CurrentUser: TempVFS.CurrentUser,
		}
	}

	usage := GetUsage()
	commands := GetCommands(vfs, usage)
	vfs.CommandMap = commands
	inputs(vfs, commands)
}
