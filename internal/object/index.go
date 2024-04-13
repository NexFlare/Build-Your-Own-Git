package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GetFileData(filesha string) []byte {
	fileContent, err := os.ReadFile(fmt.Sprintf(".git/objects/%s/%s", filesha[:2], filesha[2:]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	bytesReader := bytes.NewReader(fileContent)
	bytesWriter := &bytes.Buffer{}
	reader, err := zlib.NewReader(bytesReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong: %s\n", err.Error())
		os.Exit(1)
	}
	defer reader.Close()
	_, err = io.Copy(bytesWriter, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	return bytesWriter.Bytes()
}

func GetHeader(data []byte) ObjectHeader {
	header := ObjectHeader{}
	_, err := fmt.Sscanf(string(data), "%s %d", &header.Type, &header.Size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading header: %s\n", err.Error())
		os.Exit(1)
	}
	return header
}

func GetTreeObject(data []byte) TreeObject {
	treeObject := TreeObject{}
	treeObject.Header = GetHeader(data)
	data = data[len(data)-treeObject.Header.Size:]
	reader := bytes.NewReader(data)
	for {
		entry := &TreeEntry{
			Sha: make([]byte, 20),
		}
		mode, err := readUntil(reader, ' ')
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

		name, err := readUntil(reader, '\x00')
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

func readUntil(reader io.Reader, delim byte) (string, error) {
	buf := make([]byte, 1)
	var data []byte
	for {
		_, err := reader.Read(buf)
		if err != nil {
			return "", err
		}
		if buf[0] == delim {
			break
		}
		data = append(data, buf[0])
	}
	return string(data), nil
}

func CreateTree() string {
	return hex.EncodeToString(createTree("."))
}

func createTree(path string) []byte {
	files, err := os.ReadDir(path)
	entryList := []TreeEntry{}
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s",path,file.Name()))
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
				hashedValue:= createTree(fmt.Sprintf("%s/%s", path, file.Name()))
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
	hashedValue := getHash([]byte(treeObj))
	createObjectFile([]byte(treeObj), hex.EncodeToString(hashedValue))
	return hashedValue
}

func CreateBlog(fileContent []byte, writeToObject bool) ([]byte) {
	blobContent := fmt.Sprintf("blob %d\x00%s", len(string(fileContent)), string(fileContent))
	hashedValue := getHash([]byte(blobContent))
	if writeToObject {
		createObjectFile([]byte(blobContent), hex.EncodeToString(hashedValue))
	}
	return hashedValue
}

func createObjectFile(data []byte, hashedValue string) {
	var b bytes.Buffer
		writer := zlib.NewWriter(&b)
		writer.Write(data)
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


func getHash(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}
