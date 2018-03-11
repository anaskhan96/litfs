package filesys

import "bazil.org/fuse/fs"

type FS struct {
	RootDir *Dir
}

func (f *FS) Root() (fs.Node, error) {
	return f.RootDir, nil
}

/*func (f *FS) Destroy() {
	println("LOLLL why didn't I know about this function would've saved me so much timee")
}*/
