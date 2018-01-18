package filesys

import "bazil.org/fuse/fs"

type FS struct {
	RootDir *Dir
}

func (f *FS) Root() (fs.Node, error) {
	return f.RootDir, nil
}
