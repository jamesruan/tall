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
	Dirty bool
	Count int64
	Size  int64
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
	score    tall.HexBytes
}

func (f *FSWriteAgent) Write(data []byte) (int, error) {
	return f.w.Write(data)
}

func (f *FSWriteAgent) Close() (err error) {
	defer f.tempfile.Close()
	f.pw.Close()
	f.pr.Close()

	f.score = <-f.sch
	path := filepath.Join(f.b.entry, scoreToPath(f.score, f.b.sm.Level))
	basepath := filepath.Base(path)
	if err = os.MkdirAll(basepath, DefaultDirMode); err != nil {
		return
	}
	if err = os.Link(f.tempfile.Name(), path); err != nil {
		return
	}
	// TODO: increse Level when necessary

	return nil
}

func (f *FSWriteAgent) Score() tall.HexBytes {
	return f.score
}
