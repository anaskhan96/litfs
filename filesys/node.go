package filesys

type Node struct {
	inode uint64
	name  string
}

func NewInode(inode uint64) uint64 {
	inode += 1
	return inode
}
