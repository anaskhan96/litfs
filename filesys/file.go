package filesys

import (
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"golang.org/x/net/context"
	"sync"
)

type File struct {
	Node
	Data []byte
	sync.Mutex
}

func (file *File) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Println("Attributes for file", file.Name)
	attr.Inode = file.Inode
	attr.Mode = 0777
	attr.Size = uint64(len(file.Data))
	return nil
}

func (file *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	log.Println("Requested Read on File", file.Name)
	fuseutil.HandleRead(req, resp, file.Data)
	return nil
}

func (file *File) ReadAll(ctx context.Context) ([]byte, error) {
	log.Println("Reading all of file", file.Name)
	return []byte(file.Data), nil
}

func (file *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	log.Println("Trying to write to", file.Name, "data:", string(req.Data))
	file.Lock()
	defer file.Unlock()
	resp.Size = len(req.Data)
	file.Data = req.Data
	log.Println("Wrote to file", file.Name)
	return nil
}
func (file *File) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	log.Println("Flushing file", file.Name)
	return nil
}
func (file *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	log.Println("Open call on file", file.Name)
	return file, nil
}

func (file *File) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	log.Println("Release requested on file", file.Name)
	return nil
}

func (file *File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	return nil
}
