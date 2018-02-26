package filesys

var inode uint64

// serves as the inode table
type Node struct {
	Inode    uint64
	Name     string
	BlockNum uint32
	Size     uint64
}

func NewInode() uint64 {
	inode++
	return inode
}
