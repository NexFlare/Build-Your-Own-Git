package object

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strconv"
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
		entry.Mode, err = strconv.Atoi(mode)
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

		if _, err:=reader.Read(entry.Sha); err != nil {
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