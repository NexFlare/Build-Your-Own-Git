package handler

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/codecrafters-io/git-starter-go/internal/object"
	"github.com/codecrafters-io/git-starter-go/internal/util"
)

func InitHandler() {
	for _, dir := range []object.GitFileType{object.Root, object.Objects, object.Refs} {
		if err := os.MkdirAll(string(dir), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(string(object.Head), headFileContents, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
	}
}

func CatFileHandler() string {
	if len(os.Args) < 3 || os.Args[2] != "-p" {
		fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <sha>\n")
	}
	filesha := os.Args[3]
	fileContent := object.GetFileData(filesha)

	content := strings.Split(string(fileContent), "\x00")[1]
	return content
}

func GetHashHandler() string {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: mygit hash-object <file>\n")
		os.Exit(1)
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
	return hex.EncodeToString(hashedValue)
}

func ListTreeHandler() string {
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

	switch flag {
	case "--name-only":
		fileNames := ""
		for _, obj := range treeObject.Entries {
			if len(obj.Name) > 0 {
				fileNames = fmt.Sprintf("%s%s\n", fileNames, obj.Name)
			}
		}
		return fileNames

	default:
		return string(fileContent)
	}
}

func CreateTreeHandler() string {
	return hex.EncodeToString(object.CreateTree("."))
}


func CommitHandler() string{
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "usage: mygit commit-tree <tree_sha> -m <message>\n")
		os.Exit(1)
	}
	treeSha := os.Args[2]
	parent := ""
	message := ""
	time := time.Now()
	commitSeconds := time.Unix()
	timezone := util.GetTimeZone(time)
	authorObj := object.User {
		Email: "harshmagarwal@gmail.com",
		Name: "Harsh Agarwal",
		Seconds: commitSeconds,
		Timezone: timezone,
	}

	commiterObj := object.User {
		Email: "harshmagarwal@gmail.com",
		Name: "Harsh Agarwal",
		Seconds: commitSeconds,
		Timezone: timezone,
	}

	for i:=3; i<len(os.Args); {
		if os.Args[i] == "-m" {
			i++
			message = os.Args[i]
		} else if os.Args[i] == "-p" {
			i++
			parent = os.Args[i]
		}
		i++
	}

	commitContent := strings.Builder{}
	commitContent.WriteString(fmt.Sprintf("tree %s\n", treeSha))
	if len(parent) > 0 {
		commitContent.WriteString(fmt.Sprintf("parent %s\n", parent))
	}
	commitContent.WriteString(fmt.Sprintf("author %s <%s> %d %s\n", authorObj.Name, authorObj.Email, authorObj.Seconds, authorObj.Timezone))
	commitContent.WriteString(fmt.Sprintf("committer %s <%s> %d %s\n", commiterObj.Name, commiterObj.Email, commiterObj.Seconds, commiterObj.Timezone))
	commitContent.WriteString(fmt.Sprintf("\n%s\n", message))
	commitData := fmt.Sprintf("commit %d\x00%s", commitContent.Len(), commitContent.String())
	hashedValue := util.GetHash([]byte(commitData))
	util.CreateObjectFile([]byte(commitData), hex.EncodeToString(hashedValue))
	return hex.EncodeToString(hashedValue)
}