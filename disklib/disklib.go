package disklib

const BlockSize = 4096

import (
	"os"
)

func OpenDisk(filename string, nbytes uint64) (*os.File, error) {
	f, err := os.OpenFile("disk.bin", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return f, err
	}
	/* seek to the offset of the file */
	return f, nil
}

func ReadBlock(disk int, blocknr int) {

}

func WriteBlock(disk int, blocknr int) {

}
