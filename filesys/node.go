package filesys

var inode uint64

type Node struct {
	Inode uint64
	Name  string
}

func NewInode() uint64 {
	inode++
	return inode
}
