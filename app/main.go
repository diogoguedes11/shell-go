package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chzyer/readline"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var (
	tabCount int
	matches []string
	
)
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
type ShellCompleter struct{}

func (c *ShellCompleter) Do(line []rune, pos int) ([][]rune, int) {
    input := string(line[:pos])
    matches := findExecutables(input)
    if len(matches) == 0 {
		print("\a")
        return nil, 0
    }

    if len(matches) == 1 {
        completion := matches[0][len(input):]
	   if input + completion == matches[0] {
        	return [][]rune{[]rune(completion + " ")}, len(input)
	   }
    }
	// if len(matches) > 1 {
	// 	fmt.Printf("\n")
	// 	for _, match := range matches {
	// 		fmt.Printf("%s  ", match)
	// 	}
	// 	fmt.Printf("\n$ %s", string(line))
	// }
	print("\a")
    commonPrefix := findCommonPrefix(matches)
    if len(commonPrefix) > len(input) {
        completion := commonPrefix[len(input):]
        return [][]rune{[]rune(completion)}, len(input)
    }
    return nil, 0
}

func findExecutables(prefix string) []string {
    var matches []string
    paths := strings.Split(os.Getenv("PATH"), ":")
    
    for _, path := range paths {
        entries, err := os.ReadDir(path)
        if err != nil {
            continue
        }
        for _, e := range entries {
            name := e.Name()
            if strings.HasPrefix(name, prefix) {
                matches = append(matches, name)
            }
        }
    }
    sort.Strings(matches)
    return matches
}
// findCommonPrefix finds the longest common prefix of all strings in the slice
func findCommonPrefix(strs []string) string {
    if len(strs) == 0 {
        return ""
    }
    if len(strs) == 1 {
        return strs[0]
    }
    
    // Find the shortest string
    minLen := len(strs[0])
    for _, s := range strs {
        if len(s) < minLen {
            minLen = len(s)
        }
    }
    
    // Find the common prefix
    for i := 0; i < minLen; i++ {
        char := strs[0][i]
        for _, s := range strs {
            if s[i] != char {
                return strs[0][:i]
            }
        }
    }
    return strs[0][:minLen]
}

func main() {
	paths := strings.Split(os.Getenv("PATH"), ":")
	found := false
	config := &readline.Config{
		Prompt:       "$ ",
		AutoComplete: &ShellCompleter{},
		DisableAutoSaveHistory: false,
		EOFPrompt: "exit",
		InterruptPrompt: "^C",
	}
	rl, err := readline.NewEx(config)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error creating readline instance: ", err)
		os.Exit(1)
	}
	for {
		// fmt.Fprint(os.Stdout, "$ ")
		
		command, err := rl.Readline()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error reading command: ", err)
			os.Exit(1)
			break
		}
		trimmed := strings.TrimSpace(command)
		switch {
		case trimmed == "exit 0":
			os.Exit(0)
		case strings.Contains(trimmed, "2>>"):
			var parts []string
			parts = strings.SplitN(trimmed, "2>>", 2)
			cmdStr := strings.TrimSpace(parts[0])
			outputFile := strings.TrimSpace(parts[1])
			cmd := exec.Command("sh", "-c", cmdStr)
			dir := filepath.Dir(outputFile)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				continue
			}
			outFile, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				continue
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = outFile
			cmd.Run()
			outFile.Close()
			continue
		case strings.Contains(trimmed, "1>>"):
			var parts []string
			parts = strings.SplitN(trimmed, "1>>", 2)
			cmdStr := strings.TrimSpace(parts[0])
			outputFile := strings.TrimSpace(parts[1])

			cmd := exec.Command("sh", "-c", cmdStr)
			dir := filepath.Dir(outputFile)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				continue
			}
			f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				continue
			}
			cmd.Stdout = f
			cmd.Stderr = os.Stderr
			cmd.Run()
			f.Close()
			continue
		case strings.Contains(trimmed, ">>"):
			var parts []string
			parts = strings.SplitN(trimmed, ">>", 2)
			cmdStr := strings.TrimSpace(parts[0])
			outputFile := strings.TrimSpace(parts[1])

			cmd := exec.Command("sh", "-c", cmdStr)
			dir := filepath.Dir(outputFile)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				continue
			}
			f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				continue
			}
			cmd.Stdout = f
			cmd.Stderr = os.Stderr
			cmd.Run()
			f.Close()
			continue

		case strings.Contains(trimmed, "2>"):
			var parts []string
			parts = strings.SplitN(trimmed, "2>", 2)
			cmdStr := strings.TrimSpace(parts[0])
			outputFile := strings.TrimSpace(parts[1])
			cmd := exec.Command("sh", "-c", cmdStr)
			dir := filepath.Dir(outputFile)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				continue
			}
			outFile, err := os.Create(outputFile)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				continue
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = outFile
			cmd.Run()
			outFile.Close()
			continue
		case strings.Contains(trimmed, "1>"):
			var parts []string
			if strings.Contains(trimmed, "1>") {
				parts = strings.SplitN(trimmed, "1>", 2)
			} else {
				parts = strings.SplitN(trimmed, ">", 2)
			}
			cmdStr := strings.TrimSpace(parts[0])
			outputFile := strings.TrimSpace(parts[1])
			dir := filepath.Dir(outputFile)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				continue
			}
			cmd := exec.Command("sh", "-c", cmdStr)
			outFile, err := os.Create(outputFile)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				continue
			}

			cmd.Stdout = outFile
			cmd.Stderr = os.Stderr
			cmd.Run()
			outFile.Close()
			continue
		case strings.Contains(trimmed, ">"):
			var parts []string
			parts = strings.SplitN(trimmed, ">", 2)
			cmdStr := strings.TrimSpace(parts[0])
			outputFile := strings.TrimSpace(parts[1])

			cmd := exec.Command("sh", "-c", cmdStr)
			outFile, err := os.Create(outputFile)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				continue
			}
			defer outFile.Close()
			cmd.Stdout = outFile
			cmd.Stderr = os.Stderr
			cmd.Run()
			continue
		case strings.HasPrefix(trimmed, "cd"):
			dirPath := strings.TrimSpace(trimmed[len("cd"):])
			if dirPath == "" || dirPath == "~" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					fmt.Printf("Error while using the command cd: %v", err)
					continue
				}
				dirPath = homeDir
			}
			err := os.Chdir(dirPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: %v: No such file or directory\n", dirPath)
				continue
			}
		case strings.HasPrefix(trimmed, "echo"):
			fmt.Println(trimmed[len("echo")+1:])
		case strings.HasPrefix(trimmed, "pwd"):
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
			} else {
				fmt.Println(pwd)
			}
		case strings.HasPrefix(trimmed, "type"):
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
							log.Fatalf("Error executing the program: %s %v", programName, arguments)
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
