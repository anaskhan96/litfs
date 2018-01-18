package filesys

import (
	"log"

	"bazil.org/fuse"
	"golang.org/x/net/context"
)

type File struct {
	Node
	Data []byte
}

func (file *File) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Println("Attributes for file", file.Name)
	attr.Inode = file.Inode
	attr.Mode = 0777
	attr.Size = uint64(len(file.Data))
	return nil
}
