package filesys

type File struct {
	node Node
	data []byte
}
