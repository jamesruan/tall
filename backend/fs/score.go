package fs

import (
	"bytes"
	"crypto/sha256"
	"github.com/jamesruan/tall"
	"io"
)

func Score(d []byte) tall.HexBytes {
	r := bytes.NewBuffer(d)
	return <-ScoreFrom(r)
}

func ScoreFrom(r io.Reader) <-chan tall.HexBytes {
	ch := make(chan tall.HexBytes)
	go func() {
		defer close(ch)
		hash := sha256.New()
		io.Copy(hash, r)
		hexstring := new(tall.HexBytes)
		hexstring.FromBytes(hash.Sum(nil))
		ch <- (*hexstring)
	}()
	return ch
}
