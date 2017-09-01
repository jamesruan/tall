package fs

import (
	"github.com/jamesruan/tall"
	"os"
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
