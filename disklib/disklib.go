/* First block - metadata and structure of the filesystem
Second block - bit array of free and used blocks */

package disklib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Workiva/go-datastructures/bitarray"
	"os"
)

const (
	DISKSIZE int = 3
	BLKSIZE  int = 4096
)

type MetaBlock struct {
	Bitmap     []byte
	LowestFree int
}

var metaBlockMem MetaBlock

func OpenDisk(filename string, mbytes int) (*os.File, error) {
	if _, err := os.Stat(filename); err == nil {
		if f, err := os.OpenFile(filename, os.O_RDWR, 0666); err == nil {
			return f, nil
		}
	}
	size := uint64(mbytes * 1024 * 1024)
	fd, _ := os.Create(filename)
	fd.Seek(int64(size-1), 0)
	fd.Write([]byte{0})
	fd.Seek(0, 0)
	initBlocks(fd, size, uint64(BLKSIZE))
	return fd, nil
}

func ReadBlock(disk *os.File, blocknr int, data *[]byte) (int, error) {
	if _, err := disk.Seek(int64(blocknr*BLKSIZE), 0); err != nil {
		return 0, err
	}
	nbytes, err := disk.Read(*data)
	if err != nil {
		return 0, err
	}
	*data = bytes.Trim(*data, string(byte(0)))
	return nbytes, nil
}

func WriteBlock(disk *os.File, blocknr int, data *[]byte) (int, error) {
	zeros := make([]byte, BLKSIZE)
	disk.Write(zeros)
	if _, err := disk.Seek(int64(blocknr*BLKSIZE), 0); err != nil {
		return 0, err
	}
	nbytes, err := disk.Write(*data)
	if err != nil {
		return 0, err
	}
	return nbytes, nil
}

func initBlocks(fd *os.File, size, blksize uint64) {
	ba := bitarray.NewBitArray(size / blksize)
	ba.SetBit(0)
	ba.SetBit(1)
	data, _ := bitarray.Marshal(ba)
	metaBlockMem = MetaBlock{data, 2}
}

func MetaToDisk(f *os.File) {
	metablock, _ := json.Marshal(metaBlockMem)
	WriteBlock(f, 2, &metablock)
}

func DiskToMeta(data []byte) {
	json.Unmarshal(data, &metaBlockMem)
}