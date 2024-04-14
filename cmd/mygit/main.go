package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/internal/handler"
)

// Create a alias using alias mygit="your_git.sh"
// Usage: mygit <command> <arg1> <arg2> ...
func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		handler.InitHandler()
		fmt.Println("Initialized git directory")

	case "cat-file":
		content := handler.CatFileHandler()
		fmt.Print(content)

	case "hash-object":
		fmt.Print(handler.GetHashHandler())
	case "ls-tree":
		fmt.Print(handler.ListTreeHandler())
	
	case "write-tree":
		hashedValue := handler.CreateTreeHandler()
		fmt.Print(hashedValue)
	
	case "commit-tree":
		fmt.Print(handler.CommitHandler())

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
