package main

import (
	"fmt"
	"slices"
	"strings"
)

var AllCommands = []string{
	"cd", "mv", "history", "roothistory", "pwd", "rm", "ls",
	"fill", "mkdir", "touch", "echo", "cat", "remPerms", "whoami",
	"addPerms", "nvim", "clear", "vsh", "hostname", "sethost", "print",
	"time",
}

func (vfs *VFS) getCommandArray(name string) []string {
	fmt.Println(strings.Split(vfs.CurrentDir.Files[name].Content, ";"))
	return strings.Split(vfs.CurrentDir.Files[name].Content, ";")
}

func (vfs *VFS) vsh(name string) {
	file, exists := vfs.CurrentDir.Files[name]
	if !exists {
		fmt.Println("File ", name, "Dose not exist")
		return
	}
	if !file.Executable {
		fmt.Println("File", file.Name, "Dose not have Executable permissions")
		return
	}
	commandArray := vfs.getCommandArray(name)
	keyArray := make([][]string, len(commandArray))
	for index, value := range commandArray {
		_ = value
		/**fmt.Printf("Index: %d\nValue: %s\n", index, value) **/
		keyArray[index] = strings.Fields(commandArray[index])

		if slices.Contains(AllCommands, keyArray[index][0]) {
			fmt.Println(keyArray[index][0], ": KEYWORD")
		} else if keyArray[index][0] == "var" {
			fmt.Println(keyArray[index][0], "VAR")
		} else {
			fmt.Println(keyArray[index][0], "VALUE")
		}

	}
	// for index, value := range keyArray {
	// 	fmt.Printf("Indeex: %d\nValuee: %s\n", index, value)
	// }

	/**parts := strings.Split(file.Name, ".")
	if parts[len(parts)-1] == "vsh" {
		vfs.executeArray(vfs.getCommandArray(name))
	}
	**/

}
