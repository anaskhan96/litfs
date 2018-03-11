package filesys

import (
	"log"
	"os"
	"sync"

	"golang.org/x/net/context"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/anaskhan96/litfs/disklib"
	"syscall"
)

type Dir struct {
	Node
	Files       *[]*File
	Directories *[]*Dir
	sync.Mutex
}

func (dir *Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Println("Attributes for directory", dir.Name)
	attr.Inode = dir.Inode
	attr.Mode = os.ModeDir | 0444
	attr.BlockSize = uint32(disklib.BLKSIZE)
	return nil
}

func (dir *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("Lookup for ", name)
	if dir.Files != nil {
		for _, file := range *dir.Files {
			if file.Name == name {
				log.Println("Found match for directory lookup with size", file.Size)
				return file, nil
			}
		}
	}
	if dir.Directories != nil {
		for _, directory := range *dir.Directories {
			if directory.Name == name {
				log.Println("Found match for directory lookup")
				return directory, nil
			}
		}
	}
	return nil, fuse.ENOENT
}

func (dir *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	log.Println("Mkdir request for name", req.Name)
	newDir := &Dir{Node: Node{Name: req.Name, Inode: NewInode()}}
	directories := []*Dir{newDir}
	if dir.Directories != nil {
		directories = append(*dir.Directories, directories...)
	}
	dir.Directories = &directories
	return newDir, nil

}

func (dir *Dir) ReadDir(ctx context.Context, name string) (fs.Node, error) {
	if dir.Files != nil {
		for _, file := range *dir.Files {
			if file.Name == name {
				return file, nil
			}
		}
	}
	if dir.Directories != nil {
		for _, directory := range *dir.Directories {
			if directory.Name == name {
				return directory, nil
			}
		}
	}
	return nil, fuse.ENOENT
}

func (dir *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	log.Println("Create request for name", req.Name)
	newFile := &File{Node: Node{Name: req.Name, Inode: NewInode()}}
	files := []*File{newFile}
	if dir.Files != nil {
		files = append(files, *dir.Files...)
	}
	dir.Files = &files
	return newFile, newFile, nil
}

func (dir *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	log.Println("Remove request for ", req.Name)
	if req.Dir && dir.Directories != nil {
		newDirs := []*Dir{}
		for _, directory := range *dir.Directories {
			if directory.Name != req.Name {
				newDirs = append(newDirs, directory)
			} else {
				if directory.Files != nil {
					return fuse.Errno(syscall.ENOTEMPTY)
				}
			}
		}
		dir.Directories = &newDirs
		return nil
	} else if !req.Dir && *dir.Files != nil {
		newFiles := []*File{}
		for _, file := range *dir.Files {
			if file.Name != req.Name {
				newFiles = append(newFiles, file)
			} else {
				/* Clear the allocated blocks */
				data := make([]byte, 0)
				f, _ := disklib.OpenDisk("disklib/sda", disklib.DISKSIZE)
				defer f.Close()
				for _, i := range file.Blocks {
					disklib.WriteBlock(f, i, data)
				}
			}
		}
		dir.Files = &newFiles
		return nil
	}
	return fuse.ENOENT
}

func (dir *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("Reading all dirs")
	var children []fuse.Dirent
	if dir.Files != nil {
		for _, file := range *dir.Files {
			children = append(children, fuse.Dirent{Inode: file.Inode, Type: fuse.DT_File, Name: file.Name})
		}
	}
	if dir.Directories != nil {
		for _, directory := range *dir.Directories {
			children = append(children, fuse.Dirent{Inode: directory.Inode, Type: fuse.DT_Dir, Name: directory.Name})
		}
		log.Println(len(children), " children for dir", dir.Name)
	}
	return children, nil
}

func (dir *Dir) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
	log.Println("Rename requested from", req.OldName, "to", req.NewName)
	nd := newDir.(*Dir)
	if dir.Inode == nd.Inode {
		dir.Lock()
		defer dir.Unlock()
		for _, file := range *dir.Files {
			if file.Name == req.OldName {
				file.Name = req.NewName
				return nil
			}
		}
	} else if dir.Inode < nd.Inode {
		dir.Lock()
		defer dir.Unlock()
		nd.Lock()
		defer nd.Unlock()
		var fileDirty *File
		for _, file := range *dir.Files {
			if file.Name == req.OldName {
				fileDirty = file
				break
			}
		}
		if fileDirty != nil {
			dirExists := false
			for _, directory := range *dir.Directories {
				if directory.Name == nd.Name {
					dirExists = true
					break
				}
			}
			if dirExists {
				// Removing the file
				files := []*File{}
				for _, file := range *dir.Files {
					if file != fileDirty {
						files = append(files, file)
					}
				}
				dir.Files = &files
				// Adding in the new directory
				ndFiles := []*File{fileDirty}
				if nd.Files != nil {
					ndFiles = append(ndFiles, *nd.Files...)
				}
				nd.Files = &ndFiles
				return nil
			}
		}
	} else {
		nd.Lock()
		defer nd.Unlock()
		dir.Lock()
		defer dir.Unlock()
		var fileDirty *File
		for _, file := range *dir.Files {
			if file.Name == req.OldName {
				fileDirty = file
				break
			}
		}
		if fileDirty != nil {
			dirExists := false
			for _, directory := range *nd.Directories {
				if directory == dir {
					dirExists = true
					break
				}
			}
			if dirExists {
				// Removing the file
				files := []*File{}
				for _, file := range *dir.Files {
					if file != fileDirty {
						files = append(files, file)
					}
				}
				dir.Files = &files
				// Adding in the new directory
				ndFiles := []*File{fileDirty}
				if nd.Files != nil {
					ndFiles = append(ndFiles, *nd.Files...)
				}
				nd.Files = &ndFiles
			}
			return nil
		}
	}
	return fuse.ENOENT
}
