package main

import (
	"./disklib"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	var fsys *filesys.FS
	if _, err := os.Stat("disklib/sda"); err == nil {
		f, _ := disklib.OpenDisk("disklib/sda")
		metadataBytes := make([]byte, 2*4096)
		f.Read(metadataBytes)
		metadataBytes = bytes.Trim(metadataBytes, string(byte(0)))
		metadataMap := make(map[string]interface{})
		json.Unmarshal(metadataBytes, &metadataMap)
		rootDir, _ := metadataMap["RootDir"].(map[string]interface{})
		fmt.Println(metadataMap)
		root := setupDir(rootDir) // TODO: Inode nums not being set up properly
		fsys = &filesys.FS{&root}
	} else {
		fsys = &filesys.FS{&filesys.Dir{}}
	}

	log.Println("About to serve fs")
	go cleanup(fsys)
	if err := server.Serve(fsys); err != nil {
		log.Panicln(err)
	}

	<-conn.Ready
	if err := conn.MountError; err != nil {
		log.Panicln(err)
	}
}

func setupDir(m map[string]interface{}) filesys.Dir {
	var dir filesys.Dir
	for key, value := range m {
		if key == "Inode" {
			dir.Inode, _ = value.(uint64)
		} else if key == "Name" {
			dir.Name, _ = value.(string)
		} else if key == "Files" {
			var files []*filesys.File
			allFiles, ok := value.([]interface{})
			if !ok {
				dir.Files = nil
				continue
			}
			for _, i := range allFiles {
				val, _ := i.(map[string]interface{})
				file := setupFile(val)
				files = append(files, &file)
			}
			dir.Files = &files
		} else if key == "Directories" {
			var dirs []*filesys.Dir
			allDirs, ok := value.([]interface{})
			if !ok {
				dir.Directories = nil
				continue
			}
			for _, i := range allDirs {
				val, _ := i.(map[string]interface{})
				dirToAppend := setupDir(val)
				dirs = append(dirs, &dirToAppend)
			}
			dir.Directories = &dirs
		}
	}
	return dir
}

func setupFile(m map[string]interface{}) filesys.File {
	var file filesys.File
	for key, value := range m {
		if key == "Inode" {
			file.Inode, _ = value.(uint64)
		} else if key == "Name" {
			file.Name, _ = value.(string)
		} else if key == "Data" {
			/* Doesn't work, but then again, no need for this
			in Phase 2 as we read-write from disk */
			file.Data, _ = value.([]byte)
		}
	}
	return file
}

func cleanup(fsys *filesys.FS) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	<-done
	metadata, err := json.Marshal(&fsys)
	if err != nil {
		log.Println(err)
	}
	f, err := disklib.OpenDisk("disklib/sda")
	zeros := make([]byte, 2*4096)
	f.Write(zeros)
	f.Seek(0, 0)
	f.Write(metadata)
	fmt.Println(string(metadata))

	// Unmounting the directory
	err = fuse.Unmount("data")
	if err != nil {
		log.Println(err)
	}
}
