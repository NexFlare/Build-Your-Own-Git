package util

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"time"
)

// Get the hash of the data
func GetHash(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// Read the data from the reader until the delimiter is found
func ReadUntil(reader io.Reader, delim byte) (string, error) {
	// Creates a buffer of size 1
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

func CreateObjectFile(data []byte, hashedValue string) {
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

func GetTimeZone(time time.Time) string {
    // Get the name and offset of the local time zone
    _, zoneOffset := time.Zone()
    // Calculate hours and minutes from zone offset
    sign := "+"
    if zoneOffset < 0 {
        sign = "-"
        zoneOffset = -zoneOffset // Make the offset positive for formatting
    }
    hours := zoneOffset / 3600
    minutes := (zoneOffset % 3600) / 60

    // Format as +/-HHMM
    zoneFormatted := fmt.Sprintf("%s%02d%02d", sign, hours, minutes)

    return zoneFormatted
}