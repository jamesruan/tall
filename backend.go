package tall

import (
	"log"
)

// backend should implement following interface
type Backend interface {
	Score(data []byte) (score HexBytes) // calculate score for `data`
	Store(data []byte) (score HexBytes, err error)
	Load(score HexBytes) (data []byte, err error)
	Has(score HexBytes) (ok bool, err error)
	SetLogger(logger *log.Logger, debug bool) //operations are logged if debug is true, or else only backend running errors logged.
}

// optional interface for deleting data
type BackendScavenger interface {
	Forget(score HexBytes) // Mark the data with `score` for scavenging, Has() will return false after the call.
	Scavenge()             // Force all forgetted data to be removed. It will block until done. errors write to the Logger if is set.
}
