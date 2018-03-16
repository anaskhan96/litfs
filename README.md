# litfs

This FUSE filesystem, apart from providing the normal file I/O operations, implements persistence by emulating a single binary Unix file as disk and performing read/writes on it.

### Build and Run

```bash
go get github.com/anaskhan96/litfs
cd $GOPATH/src/github.com/anaskhan96/litfs
go run main.go data # data/ is the directory on which to mount the filesystem on
```

Run `umount <path-to-directory>` to unmount the filesystem.