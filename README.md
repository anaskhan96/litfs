# litfs

This FUSE filesystem, apart from providing the normal file I/O operations, implements persistence by emulating a single binary Unix file as disk and performing read/writes on it.

## Build and Run

```bash
go get github.com/anaskhan96/litfs
cd $GOPATH/src/github.com/anaskhan96/litfs
go run main.go data # data/ is the directory on which to mount the filesystem on
```

Run `umount <path-to-directory>` to unmount the filesystem.

## File System Characteristics

- Create, remove a directory
- Create, remove, read from, and write to files inside a directory
- Copy, move the contents of a file to another, across directories

## Persistence Implementation

Keeping a block size of `4096` bytes:

- A serialized form of the tree representation of the filesystem is stored in the first block
- A structure containing two components - a bitmap indicating free and allocated blocks in the filesystem and an integer containing the lowest free block at the moment - is serialized and stored in the second block
- File data is stored from the third block onwards, with a block as a whole being allocated to/deallocated from the file