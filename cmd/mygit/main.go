package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

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
		blobContent := fmt.Sprintf("blob %d\x00%s", len(string(fileContent)), string(fileContent))
		hashedValue := getHash([]byte(blobContent))
		if writeToObj {
			var b bytes.Buffer
			writer := zlib.NewWriter(&b)
			writer.Write([]byte(blobContent))
			writer.Close()
			if err := os.MkdirAll(fmt.Sprintf(".git/objects/%s", hashedValue[:2]), 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
			err := os.WriteFile(fmt.Sprintf(".git/objects/%s/%s", hashedValue[:2], hashedValue[2:]), b.Bytes(), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err.Error())
				os.Exit(1)
			}
		}
		fmt.Print(hashedValue)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}


func getHash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	hashedValue := hex.EncodeToString(h.Sum(nil))
	return hashedValue
}