package object

type ObjectType string

const (
	Blob ObjectType = "blob"
	Tree ObjectType = "tree"
)

type ObjectHeader struct {
	Type ObjectType
	Size int
}

type TreeEntry struct {
	Mode int
	Name string
	Sha []byte
}

type TreeObject struct {
	Header ObjectHeader
	Entries []TreeEntry
}