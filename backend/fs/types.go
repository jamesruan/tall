package fs

import (
	"github.com/jamesruan/tall"
	"io"
	"os"
	"path/filepath"
)

const (
	JOURNALPATH    = "journals/"
	TEMPPATH       = "temp/"
	DATAPATH       = "data/"
	STATSPATH      = "stats/"
	SUPERMATPATH   = "supermata"
	DefaultMode    = os.FileMode(0644)
	DefaultDirMode = os.FileMode(0755)
	RootScore      = tall.HexBytes("")
)

type SuperMeta struct {
	Level int
}

func (s *SuperMeta) Store(f *os.File) (err error) {
	return StoreStructToFile(s, f)
}

func (s *SuperMeta) Load(f *os.File) (err error) {
	return LoadStructFromFile(s, f)
}

type FSWriteAgent struct {
	b        *FSBackend
	sch      <-chan tall.HexBytes
	tempfile *os.File
	w        io.Writer
	pw       io.WriteCloser
	pr       io.ReadCloser
}

type FSReadAgent struct {
	*os.File
}

func (f *FSWriteAgent) Write(data []byte) (int, error) {
	return f.w.Write(data)
}

func (f *FSWriteAgent) Commit() (score tall.HexBytes, err error) {
	defer f.tempfile.Close()
	f.pw.Close()
	f.pr.Close()

	score = <-f.sch
	path := filepath.Join(f.b.entry, scoreToPath(score, f.b.sm.Level))
	basepath := filepath.Base(path)
	if err = os.MkdirAll(basepath, DefaultDirMode); err != nil {
		return
	}
	if err = os.Link(f.tempfile.Name(), path); err != nil {
		return
	}
	// TODO: increse Level when necessary

	return
}
