package tall

// implements FileSystem
// Empty dir is treated as "."
type Dir struct {
	Path string
	be   Backend
}

func (d Dir) Open(name string) (File, error) {
	//TODO
	return nil, nil
}

func NewDir(backend Backend, path string) {
}
