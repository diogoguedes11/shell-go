package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage
	paths := strings.Split(os.Getenv("PATH"), ":")
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprint(os.Stderr, "Error reading command: ", err)
			os.Exit(1)
			break
		}
		trimmed := strings.TrimSpace(command)
		switch {
			case trimmed == "exit 0":
				os.Exit(0)
			case strings.HasPrefix(trimmed,"echo"):
				fmt.Println(trimmed[len("echo")+1:])
			case strings.HasPrefix(trimmed,"type"):
				found := false
				cmdName := trimmed[len("type")+1:]
				for _, path := range paths {
					fullPath := path + "/" + cmdName
					if fileInfo, err := os.Stat(fullPath); err == nil {
						if fileInfo.Mode().IsRegular() && fileInfo.Mode()&0111 != 0 {
							fmt.Printf("%s is %s\n", cmdName, fullPath)
							found = true
							break
						}
					}
				}
				if !found {
					fmt.Println(trimmed[len("type")+1:] + ": not found")
				}
			default:
				fmt.Println(trimmed + ": command not found")
		}
		
		
	}
}
