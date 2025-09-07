package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	paths := strings.Split(os.Getenv("PATH"), ":")
	found := false
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
			case strings.Contains(trimmed,"1>") || strings.Contains(trimmed,">"):
				var parts []string;
				if strings.Contains(trimmed,"1>") {
					parts = strings.SplitN(trimmed,"1>",2)
				} else {
					parts = strings.SplitN(trimmed,">",2)
				}
				cmdStr := strings.TrimSpace(parts[0])
				outputFile := strings.TrimSpace(parts[1])

				cmd := exec.Command("sh","-c",cmdStr)
				outFile , err := os.Create(outputFile)

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
					continue
				}
				defer outFile.Close()
				cmd.Stdout = outFile
				cmd.Stderr = os.Stderr
				cmd.Run()
				continue
			case strings.Contains(trimmed,"2>"):
				var parts []string;
				parts = strings.SplitN(trimmed,"2>",2)
				cmdStr := strings.TrimSpace(parts[0])
				outputFile := strings.TrimSpace(parts[1])

				cmd := exec.Command("sh","-c",cmdStr)
				outFile , err := os.Create(outputFile)

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
					continue
				}
				defer outFile.Close()
				cmd.Stderr = outFile
				cmd.Stdout = os.Stdout
				cmd.Run()
				continue
					
			case strings.HasPrefix(trimmed,"echo"):
				fmt.Println(trimmed[len("echo")+1:])
			case strings.HasPrefix(trimmed,"cd"):
				fullPath := trimmed[len("cd")+1:]
				if _, err := os.ReadDir(fullPath); err != nil  && fullPath != "~" {
					fmt.Fprintf(os.Stderr, "cd: %v: No such file or directory\n", fullPath)
				} else if fullPath == "~" {
					path := os.Getenv("HOME")
					os.Chdir(path)
				}else {
					os.Chdir(fullPath)
				}
				
			case strings.HasPrefix(trimmed,"pwd"):
				pwd,err := os.Getwd()
				if err != nil {
					fmt.Fprintf(os.Stderr,"pwd: %v\n",err)
				}else {
					fmt.Println(pwd)
				}
			case strings.HasPrefix(trimmed,"type"):
				cmdName := trimmed[len("type")+1:]
				if cmdName == "echo" || cmdName == "type" || cmdName == "exit" || cmdName == "pwd" {
					fmt.Printf("%s is a shell builtin\n", cmdName)
					break
				} else {
					found = false
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
				}
				if !found {
					fmt.Println(cmdName + ": not found")
				}
			default:
				parts := strings.Fields(trimmed)
				programName := parts[0]
				arguments := parts[1:]
				found = false
				for _, path := range paths {
					fullPath := path + "/" + programName
					if fileInfo, err := os.Stat(fullPath); err == nil {
						if fileInfo.Mode().IsRegular() && fileInfo.Mode()&0111 != 0 {
							cmd := exec.Command(programName, arguments...)
							cmd.Stdout = os.Stdout // allows me to get the output in my shell
							cmd.Stderr = os.Stderr // allows me to get the error output in my shell
							err := cmd.Run()
							if err != nil {
								log.Fatalf("Error executing the program: %s %v",programName,arguments)
								return
							}
							found = true
							break
						} 
					}
				}
				if !found {
					fmt.Println(programName + ": not found")
				}
		}
		
	}
}
