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
func findLongestCommonPrefix(input string,matches []string) string {
	if len(matches) == 0{
		return ""
	}
	minLen := -1
	for _, m := range matches {
		if minLen == -1 || len(m) < minLen {
			minLen = len(m)
		}
	}
	if minLen == -1 {
		return ""
	}
	for i := 0; i < minLen; i++ {
    		ch := matches[0][i]
		for _, m := range matches {
			if m[i] != ch {
				if i < len(input) {
					print("\a")
					return ""
				}
					print("\a")
					return matches[0][:i]	
			} 
		}
	}
	if minLen < len(input) {
		return ""
	}
	// If we reach here, it means we found a common prefix
	return matches[0][:minLen]
}
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
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}
	for _, b := range builtins {
		if strings.HasPrefix(b, input) && !contains(matches,b) {
			matches = append(matches, b)
		}
	}
	sort.Strings(matches)
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
	if len(matches) > 1 {
        commonPrefix := findLongestCommonPrefix(input, matches)
        // only take the tail if the common prefix is longer than the input
        if len(commonPrefix) > len(input) {
            completion := commonPrefix[len(input):]
            return [][]rune{[]rune(completion)}, len(input)
        }
        fmt.Fprint(os.Stdout, "\n")
        for i, m := range matches {
            if i > 0 {
                fmt.Fprint(os.Stdout, "  ")
            }
            fmt.Fprint(os.Stdout, m)
        }
        fmt.Fprintf(os.Stdout, "\n$ %s", input)
        return nil, 0
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
            if strings.HasPrefix(name, prefix) && !contains(matches,name) {
                matches = append(matches, name)
            }
        }
    }
    sort.Strings(matches)
    return matches
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
			arg := strings.TrimSpace(strings.TrimPrefix(trimmed, "echo"))
			if len(arg) >= 2 && ((strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'")) || (strings.HasPrefix(arg, `"` ) && strings.HasSuffix(arg, `"`))) {
				fmt.Fprintln(os.Stdout,arg[1 :len(arg)-1])
			} else {
				fmt.Fprintln(os.Stdout,  strings.Join(strings.Fields(arg), " "))
			}
		case strings.HasPrefix(trimmed,"exit"):
			os.Exit(0)
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
				// if ((strings.HasPrefix(path, "'") && strings.HasSuffix(path, "'")) || (strings.HasPrefix(path, `"` ) && strings.HasSuffix(path, `"`))) {
				// 	fmt.Fprintln(os.Stdout,fullPath[1:len(fullPath)-1])
				// } 
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
