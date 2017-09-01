package tall

import (
	"io"
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
}

// optional interface for deleting data
type BackendScavenger interface {
	Forget(score HexBytes) (err error) // Mark the data with `score` for scavenging, Has() will return false after the call.
	Scavenge() (err error)             // Force all forgetted data to be removed. It will block until done.
}
