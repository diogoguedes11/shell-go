package main

import (
	"bufio"
	"fmt"
	"os"
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
		if command[:len(command)-1] == "exit 0" {
			os.Exit(0)
			break
		}
		fmt.Println(command[:len(command)-1] + ": command not found")
		
	}
}
