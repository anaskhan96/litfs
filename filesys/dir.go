package filesys

type Dir struct {
	node  Node
	files *[]*File
	dirs  *[]*Dir
}
