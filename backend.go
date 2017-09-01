package tall

import (
	"io"
	"log"
)

type Reader interface {
	io.ReadSeeker
	io.ReaderAt
	io.Closer
	Score() (score HexBytes)
}

type Writer interface {
	io.Writer
	Commit() (score HexBytes, err error)
}

// backend should implement following interface
type Backend interface {
	Score(data []byte) (score HexBytes)      // calculate score for `data`
	Open(score HexBytes) (Reader, err error) // open a storage for read, call Close() to close it.
	Create() (Writer, err error)             // create a new storage to write, call Commit() to finalize it.
	Has(score HexBytes) (ok bool, err error)
	SetLogger(logger *log.Logger, debug bool) //operations are logged if debug is true, or else only backend running errors logged.
}

// optional interface for deleting data
type BackendScavenger interface {
	Forget(score HexBytes) // Mark the data with `score` for scavenging, Has() will return false after the call.
	Scavenge()             // Force all forgetted data to be removed. It will block until done. errors write to the Logger if is set.
}
