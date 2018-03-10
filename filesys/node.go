package filesys

var inode uint64

// serves as the inode table
type Node struct {
	Inode  uint64
	Name   string
	Blocks []int
	Size   uint64
}

func NewInode() uint64 {
	inode++
	return inode
}

func InitInode(n uint64) {
	inode = n
}
