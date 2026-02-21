package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type Action string

const (
	Add Action = "add"
	Remove Action = "remove"
)

const path string = "C:\\Windows\\System32\\drivers\\etc\\hosts"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No Arguments Provided")
		return
	}

	var action Action
	var domain string
	var ip string

	action = Action(args[0])

	if action != Add && action != Remove {
		fmt.Println("Invalid action. Use 'add' or 'remove'")
		return
	} 

	for i := 1; i < len(args); i++ {
		if args[i] == "ip" && i+1 < len(args) {
			ip = args[i+1]
			i++
		} else if args[i] == "domain" && i+1 < len(args) {
			domain = args[i+1]
			i++
		} 
	}

	if ip == "" || domain == "" {
		fmt.Println("Both ip and domain are required")
		return
	}

	target := ip + " " + domain

	if action == Add {
		file, err := os.OpenFile(path, os.O_APPEND | os.O_WRONLY, 0644,)
		
		if err != nil {
			panic(err)
		}
		
		defer file.Close()

		_, err = file.WriteString(target + "\r\n")
		if err != nil {
			panic(err)
		}

		fmt.Println("Entry added successfully")
		return

	} else if action == Remove {
		
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}

		tempPath := filepath.Join(filepath.Dir(path), "hosts.tmp")
		tempFile, err := os.Create(tempPath)
		if err != nil {
			file.Close()
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		found := false

		for scanner.Scan(){
			line := scanner.Text()

			if line == target {
				found = true
				continue
			}

			_, err := tempFile.WriteString(line + "\r\n")
			if err != nil {
				file.Close()
				tempFile.Close()
				panic(err)
			}		
		}

		if err := scanner.Err(); err != nil {
			file.Close()
			tempFile.Close()
			panic(err)
		}

		file.Close()
		tempFile.Close()

		// Replace original file
		err = os.Remove(path)
		if err != nil {
			os.Remove(tempPath)
			panic(err)
		}

		err = os.Rename(tempPath, path)
		if err != nil {
			panic(err)
		}

		if found {
			fmt.Println("Entry removed successfully")
		} else {
			fmt.Println("Entry not found")
		}
	}
}