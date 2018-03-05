package disklib

import (
	"os"
)

const (
	DISKSIZE int = 3
	BLKSIZE  int = 4096
)

func OpenDisk(filename string, mbytes int) (*os.File, error) {
	if _, err := os.Stat(filename); err == nil {
		if f, err := os.OpenFile(filename, os.O_RDWR, 0666); err == nil {
			return f, nil
		}
	}
	size := int64(mbytes * 1024 * 1024) // Creating a 5 MB disk file
	fd, _ := os.Create(filename)
	fd.Seek(size-1, 0)
	fd.Write([]byte{0})
	fd.Seek(0, 0)
	return fd, nil
}

func ReadBlock(disk *os.File, blocknr int, data *[]byte) {
	disk.Seek(int64(blocknr*BLKSIZE), 0)
	disk.Read(*data)
}

func WriteBlock(disk *os.File, blocknr int, data *[]byte) {
	zeros := make([]byte, BLKSIZE)
	disk.Write(zeros)
	disk.Seek(int64(blocknr*BLKSIZE), 0)
	disk.Write(*data)
}
