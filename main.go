package main

import (
	"fmt"
	"log"
	"os"

	"./filesys"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Provide directory for mounting")
		os.Exit(1)
	}

	mountpoint := os.Args[1]

	conn, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := fs.New(conn, nil)

	fsys := &filesys.FS{
		&filesys.Dir{Node: filesys.Node{Name: "head", Inode: filesys.NewInode()}, Files: &[]*filesys.File{
			&filesys.File{Node: filesys.Node{Name: "hello", Inode: filesys.NewInode()}, Data: []byte("hello world!")},
			&filesys.File{Node: filesys.Node{Name: "aybbg", Inode: filesys.NewInode()}, Data: []byte("send notes")},
		}, Directories: &[]*filesys.Dir{
			&filesys.Dir{Node: filesys.Node{Name: "left", Inode: filesys.NewInode()}, Files: &[]*filesys.File{
				&filesys.File{Node: filesys.Node{Name: "yo", Inode: filesys.NewInode()}, Data: []byte("ayy lmao")},
			},
			},
			&filesys.Dir{Node: filesys.Node{Name: "right", Inode: filesys.NewInode()}, Files: &[]*filesys.File{
				&filesys.File{Node: filesys.Node{Name: "hey", Inode: filesys.NewInode()}, Data: []byte("heeey, thats pretty good")},
			},
			},
		},
		},
	}
	log.Println("About to serve fs")
	if err := server.Serve(fsys); err != nil {
		log.Panicln(err)
	}

	<-conn.Ready
	if err := conn.MountError; err != nil {
		log.Panicln(err)
	}
}
