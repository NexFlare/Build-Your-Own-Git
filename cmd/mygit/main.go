package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"os"

	"github.com/codecrafters-io/git-starter-go/internal/object"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		if len(os.Args) < 3 || os.Args[2] != "-p" {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <sha>\n")
		}
		filesha := os.Args[3]
		fileContent := object.GetFileData(filesha)
		// fmt.Println("Reading file", filesha)

		content := strings.Split(string(fileContent), "\x00")[1]
		fmt.Print(content)

	case "hash-object":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: mygit hash-object <file>\n")
		}
		writeToObj := false
		fileName := os.Args[2]
		if len(os.Args) > 3 && os.Args[2] == "-w" {
			writeToObj = true
			fileName = os.Args[3]
		}
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
			os.Exit(1)
		}

		hashedValue := object.CreateBlog(fileContent, writeToObj)
		fmt.Print(hex.EncodeToString(hashedValue))
	case "ls-tree":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: mygit ls-tree <sha>\n")
		}
		flag := ""
		if len(os.Args) > 3 {
			flag = os.Args[2]
		}
		filesha := os.Args[3]
		fileContent := object.GetFileData(filesha)
		treeObject := object.GetTreeObject(fileContent)
		// fileContent := getFileContent(filesha)
		switch flag {
		case "--name-only":
			fileNames := ""
			for _, obj := range treeObject.Entries {
				if len(obj.Name) > 0 {
					fileNames = fmt.Sprintf("%s%s\n", fileNames, obj.Name)
				}
			}
			fmt.Print(fileNames)

		default:
			fmt.Print(string(fileContent))
		}
	
	case "write-tree":
		hashedValue := object.CreateTree()
		fmt.Print(hashedValue)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
