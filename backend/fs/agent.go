package fs

import (
	"github.com/jamesruan/tall"
	"io"
	"os"
	"path/filepath"
)

type FSWriteAgent struct {
	b        *FSBackend
	sch      <-chan tall.HexBytes
	tempfile *os.File
	w        io.Writer
	pw       io.WriteCloser
	pr       io.ReadCloser
}

// implements tall.Reader
type FSReadAgent struct {
	*os.File
}

func (f *FSReadAgent) Score() tall.HexBytes {
	b := new(tall.HexBytes)
	name := filepath.Base(f.File.Name())
	if err := b.FromString(name); err != nil {
		panic(err)
	} else {
		return *b
	}
}

// implements tall.Writer
func (f *FSWriteAgent) Write(data []byte) (int, error) {
	return f.w.Write(data)
}

func (f *FSWriteAgent) Commit() (score tall.HexBytes, err error) {
	defer f.tempfile.Close()
	f.pw.Close()
	f.pr.Close()

	score = <-f.sch
	path := filepath.Join(f.b.entry, scoreToPath(score, f.b.sm.Level))
	dirpath := filepath.Dir(path)
	if err = os.MkdirAll(dirpath, DefaultDirMode); err != nil {
		return
	}
	if err = os.Link(f.tempfile.Name(), path); err != nil {
		return
	}
	if err = os.Remove(f.tempfile.Name()); err != nil {
		return
	}
	// TODO: increse Level when necessary

	return
}
