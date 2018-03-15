package filesys

import (
	"bazil.org/fuse/fs"
	"encoding/json"
	"github.com/anaskhan96/litfs/disklib"
	"log"
)

type FS struct {
	RootDir *Dir
}

func (f *FS) Root() (fs.Node, error) {
	return f.RootDir, nil
}

func (fsys *FS) Destroy() {
	metadata, err := json.Marshal(&fsys)
	if err != nil {
		log.Println(err)
	}
	f, err := disklib.OpenDisk("disklib/sda", disklib.DISKSIZE)
	defer f.Close()
	disklib.WriteBlock(f, 0, metadata)
	disklib.MetaToDisk(f)
}
