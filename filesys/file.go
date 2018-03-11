package filesys

import (
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"github.com/anaskhan96/litfs/disklib"
	"golang.org/x/net/context"
	"math"
	"sync"
)

type File struct {
	Node
	//Data []byte
	sync.Mutex
}

func (file *File) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Println("Attributes for file", file.Name)
	attr.Inode = file.Inode
	attr.Mode = 0777
	attr.Size = file.Size
	attr.BlockSize = uint32(disklib.BLKSIZE)
	attr.Blocks = uint64(len(file.Blocks))
	return nil
}

/* Look at this function later because it's not supposed to return the whole data */
func (file *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	log.Println("Requested Read on File", file.Name)
	data := make([]byte, 0, 10)
	f, _ := disklib.OpenDisk("disklib/sda", disklib.DISKSIZE)
	defer f.Close()
	for i := 0; i < len(file.Blocks); i++ {
		blockData := make([]byte, disklib.BLKSIZE)
		disklib.ReadBlock(f, file.Blocks[i], &blockData)
		data = append(data, blockData...)
	}
	fuseutil.HandleRead(req, resp, data)
	return nil
}

func (file *File) ReadAll(ctx context.Context) ([]byte, error) {
	log.Println("Reading all of file", file.Name)
	data := make([]byte, 0, 10)
	f, _ := disklib.OpenDisk("disklib/sda", disklib.DISKSIZE)
	defer f.Close()
	for i := 0; i < len(file.Blocks); i++ {
		blockData := make([]byte, disklib.BLKSIZE)
		disklib.ReadBlock(f, file.Blocks[i], &blockData)
		data = append(data, blockData...)
	}
	return []byte(data), nil
}

func (file *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	log.Println("Trying to write to", file.Name, "data:", string(req.Data))
	file.Lock()
	defer file.Unlock()
	resp.Size = len(req.Data)
	file.Size = uint64(len(req.Data))
	numBlocks := int(math.Ceil(float64(file.Size) / float64(disklib.BLKSIZE)))
	blocks := make([]int, numBlocks)
	f, _ := disklib.OpenDisk("disklib/sda", disklib.DISKSIZE)
	k := 0
	for i := 0; i < numBlocks; i++ {
		blocknr := disklib.GetLowestFreeBlock()
		var data []byte
		if i == numBlocks-1 {
			data = req.Data[k:]
		} else {
			data = req.Data[k:disklib.BLKSIZE]
		}
		k += disklib.BLKSIZE
		nbytes, err := disklib.WriteBlock(f, blocknr, data)
		log.Println("BYTES WRITTEN:", nbytes, "ERR:", err)
		blocks[i] = blocknr
	}
	f.Close()
	file.Blocks = blocks
	log.Println(file.Blocks)
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
	log.Println("Fsync requested on file", file.Name)
	return nil
}
