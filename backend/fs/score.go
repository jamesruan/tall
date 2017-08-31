package fs

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/jamesruan/tall"
	"io"
)

func Score(d []byte) tall.HexBytes {
	r := bytes.NewBuffer(d)
	return <-ScoreFrom(r)
}

func ScoreFrom(r io.Reader) <-chan tall.HexBytes {
	ch := make(chan tall.HexBytes, 1)
	go func() {
		hash := sha256.New()
		io.Copy(hash, r)
		hexstring := hex.EncodeToString(hash.Sum(nil))
		ch <- tall.HexBytes(hexstring)
	}()
	return ch
}
