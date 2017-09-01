package tall

import (
	"encoding/hex"
	"io"
	"os"
)

// implements http.File
type File interface {
	io.ReadSeeker
	io.ReaderAt
	io.Closer
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

type FileSystem interface {
	Open(name string) (File, error) //implements http.FileSystem
}

// hex representation of bytes
type HexBytes string

func (h HexBytes) String() string {
	return string(h)
}

func (h *HexBytes) FromString(s string) error {
	if _, err := hex.DecodeString(s); err != nil {
		return err
	} else {
		*h = HexBytes(s)
		return nil
	}
}

func (h *HexBytes) FromBytes(b []byte) {
	*h = HexBytes(hex.EncodeToString(b))
}

func (h HexBytes) Bytes() []byte {
	if b, err := hex.DecodeString(string(h)); err != nil {
		panic(err)
	} else {
		return b
	}
}

type DirInfo struct {
	Entries []*FileStat
	index   map[string]int // maps file name to index in FileStat
}
