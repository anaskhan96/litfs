package filesys

import "bazil.org/fuse/fs"

type FS struct {
	root *Dir
}

func (f *FS) Root() (fs.Node, error) {
	return f.root, nil
}
