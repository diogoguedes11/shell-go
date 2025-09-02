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
				if trimmed[len("type")+1:] == "echo" || trimmed[len("type")+1:] == "exit" || trimmed[len("type")+1:] == "type"{
					fmt.Printf("%s is a shell builtin\n", trimmed[len("type")+1:])
				} else {
					fmt.Println(trimmed[len("type")+1:] + ": not found")
				}
			default:
				fmt.Println(trimmed + ": command not found")
		}
		
		
	}
}
