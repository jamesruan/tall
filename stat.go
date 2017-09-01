package tall

import (
	"os"
	"time"
)

// implements os.FileInfo
type FileStat struct {
	Fscore   HexBytes    `json:"score"`
	Fpscore  HexBytes    `json:"pscore"` // score of previous version
	Fname    string      `json:"name"`
	Fsize    int64       `json:"size"`
	Fmode    os.FileMode `json:"mode"`
	FmodTime time.Time   `json:"time"`
}

func (fd FileStat) Name() string {
	return fd.Fname
}

func (fd FileStat) Size() int64 {
	return fd.Fsize
}

func (fd FileStat) Mode() os.FileMode {
	return fd.Fmode
}

func (fd FileStat) ModTime() time.Time {
	return fd.FmodTime
}

func (fd *FileStat) Sys() interface{} {
	return fd
}

func (fd FileStat) IsDir() bool {
	return fd.Fmode.IsDir()
}
