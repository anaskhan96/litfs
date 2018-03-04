package disklib

import (
	"os"
)

func OpenDisk(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); err == nil {
		if f, err := os.OpenFile(filename, os.O_RDWR, 0666); err == nil {
			return f, nil
		}
	}
	size := int64(5 * 1024 * 1024) // Creating a 5 MB disk file
	fd, _ := os.Create(filename)
	fd.Seek(size-1, 0)
	fd.Write([]byte{0})
	fd.Seek(0, 0)
	return fd, nil
}

func ReadBlock(disk int, blocknr int) {

}

func WriteBlock(disk int, blocknr int) {

}
