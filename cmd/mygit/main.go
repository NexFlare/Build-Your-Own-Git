package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"strings"

	// Uncomment this block to pass the first stage!

	"os"
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
		// fmt.Println("Reading file", filesha)
		fileContent, err := os.ReadFile(fmt.Sprintf(".git/objects/%s/%s", filesha[:2], filesha[2:]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
			os.Exit(1)
		}
		bytesReader := bytes.NewReader(fileContent)
		bytesWriter := &bytes.Buffer{}
		reader, err := zlib.NewReader(bytesReader)
		defer reader.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong: %s\n", err.Error())
			os.Exit(1)
		}
		_, err = io.Copy(bytesWriter, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
			os.Exit(1)
		}
		// 
		// fmt.Println(bytesWriter.String())
		_ = strings.Split(bytesWriter.String(), "")[0]
		content := strings.Split(bytesWriter.String(), "\x00")[1]
		fmt.Print(content)
	
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
