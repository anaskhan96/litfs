/* First block - metadata and structure of the filesystem
Second block - bit array of free and used blocks */

package disklib

import (
	"bytes"
	"encoding/json"
	"github.com/Workiva/go-datastructures/bitarray"
	"log"
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

func WriteBlock(disk *os.File, blocknr int, data []byte) (int, error) {
	zeros := make([]byte, BLKSIZE)
	if _, err := disk.Seek(int64(blocknr*BLKSIZE), 0); err != nil {
		return 0, err
	}
	disk.Write(zeros)
	if len(data) == 0 {
		updateBlocks("del", blocknr)
		return 0, nil
	}
	if _, err := disk.Seek(-int64(BLKSIZE), 1); err != nil {
		return 0, err
	}
	nbytes, err := disk.Write(data)
	if err != nil {
		return 0, err
	}
	updateBlocks("set", blocknr)
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
	WriteBlock(f, 1, metablock)
}

func DiskToMeta(data []byte) {
	json.Unmarshal(data, &metaBlockMem)
}

func updateBlocks(operation string, blocknr int) {
	if operation == "del" {
		ba, _ := bitarray.Unmarshal(metaBlockMem.Bitmap)
		ba.ClearBit(uint64(blocknr))
		data, _ := bitarray.Marshal(ba)
		metaBlockMem.Bitmap = data
		if blocknr < metaBlockMem.LowestFree {
			metaBlockMem.LowestFree = blocknr
		}
	} else if operation == "set" {
		ba, _ := bitarray.Unmarshal(metaBlockMem.Bitmap)
		if ok, _ := ba.GetBit(uint64(blocknr)); !ok {
			ba.SetBit(uint64(blocknr))
			data, _ := bitarray.Marshal(ba)
			metaBlockMem.Bitmap = data
			if blocknr == metaBlockMem.LowestFree {
				i := blocknr + 1
				for i < int(ba.Capacity()) {
					if ok, _ := ba.GetBit(uint64(i)); !ok {
						metaBlockMem.LowestFree = i
					}
					i++
				}
			}
		}
	}
	log.Println("New lowest free block:", metaBlockMem.LowestFree)
}

func GetLowestFreeBlock() int {
	return metaBlockMem.LowestFree
}
