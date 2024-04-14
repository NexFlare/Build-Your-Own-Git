package object

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/git-starter-go/internal/util"
)

func GetTreeObject(data []byte) TreeObject {
	treeObject := TreeObject{}
	treeObject.Header = GetHeader(data)
	data = data[len(data)-treeObject.Header.Size:]
	reader := bytes.NewReader(data)
	for {
		entry := &TreeEntry{
			Sha: make([]byte, 20),
		}
		mode, err := util.ReadUntil(reader, ' ')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error while reading file: %s\n", err.Error())
			}
			break
		}
		entry.Mode = Mode(mode)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error while reading file: %s\n", err.Error())
			}
			break
		}

		name, err := util.ReadUntil(reader, '\x00')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error while reading file: %s\n", err.Error())
			}
			break
		}
		entry.Name = name

		if _, err := reader.Read(entry.Sha); err != nil {
			break
		}
		treeObject.Entries = append(treeObject.Entries, *entry)
	}
	return treeObject
}

func CreateTree(path string) []byte {
	files, err := os.ReadDir(path)
	entryList := []TreeEntry{}
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
					os.Exit(1)
				}
				hashedValue := CreateBlog(fileContent, true)
				entryList = append(entryList, TreeEntry{
					Mode: FileMode,
					Name: file.Name(),
					Sha:  hashedValue,
				})
			} else {
				if file.Name() == ".git" {
					continue
				}
				hashedValue := CreateTree(fmt.Sprintf("%s/%s", path, file.Name()))
				entryList = append(entryList, TreeEntry{
					Mode: DirMode,
					Name: file.Name(),
					Sha:  hashedValue,
				})
			}
		}
	}
	treeContent := ""
	for _, entry := range entryList {
		treeContent = fmt.Sprintf("%s%s %s\x00%s", treeContent, entry.Mode, entry.Name, entry.Sha)
	}
	treeObj := fmt.Sprintf("tree %d\x00%s", len(treeContent), treeContent)
	hashedValue := util.GetHash([]byte(treeObj))
	util.CreateObjectFile([]byte(treeObj), hex.EncodeToString(hashedValue))
	return hashedValue
}
