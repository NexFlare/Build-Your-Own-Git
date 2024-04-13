package object

import "fmt"

type ObjectType string
type Mode string


const (
	Blob ObjectType = "blob"
	Tree ObjectType = "tree"
)

const (
	FileMode       Mode = "100644"
	DirMode        Mode = "40000"
	ExecutableMode Mode = "100755"
	SymlinkMode    Mode = "120000"
)

type ObjectHeader struct {
	Type ObjectType
	Size int
}

type TreeEntry struct {
	Mode Mode
	Name string
	Sha  []byte
}

type TreeObject struct {
	Header  ObjectHeader
	Entries []TreeEntry
}

type BlobObject struct {
	Header ObjectHeader
	Content string
}

func (b BlobObject) String() string {
	return fmt.Sprintf("%s %d\x00%s", b.Header.Type, b.Header.Size, b.Content)
}
