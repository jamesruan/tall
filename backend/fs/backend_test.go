package fs

import (
	"bytes"
	"github.com/jamesruan/tall"
	"io"
	"os"
	"testing"
)

func TestMake(t *testing.T) {
	path := "test"
	if err := Make(path, true); err != nil {
		t.Error(err)
	}
	if err := Make(path, false); err == nil {
		t.Error("exist detect\n")
	}
	if err := Make(path, true); err != nil {
		t.Error(err)
	}
	os.RemoveAll(path)
}

func TestStore(t *testing.T) {
	var err error
	TEST := []byte("test")
	path := "test"
	if err := Make(path, true); err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(path)

	var be *FSBackend

	if be, err = New(path); err != nil {
		t.Error(err)
	}

	var w tall.Writer
	if w, err = be.Create(); err != nil {
		t.Error(err)
	}

	if _, err = w.Write(TEST); err != nil {
		t.Error(err)
	}

	var score tall.HexBytes

	if score, err = w.Commit(); err != nil {
		t.Error(err)
	}

	t.Logf("write done: %s\n", score)

	if score != be.Score(TEST) {
		t.Error("bad score\n")
	}

	if ok, err := be.Has(score); err != nil || !ok {
		t.Error(err)
	}

	var r tall.Reader
	if r, err = be.Open(score); err != nil {
		t.Error(err)
	}
	defer r.Close()

	t.Logf("Score(): %s\n", r.Score())
	if r.Score().String() != score.String() {
		t.Errorf("score mismatch %s %s\n", r.Score(), score)
	}

	buf := new(bytes.Buffer)

	if _, err = io.Copy(buf, r); err != nil {
		t.Error(err)
	}

	if string(buf.Bytes()) != string(TEST) {
		t.Errorf("error read after write %s", buf.Bytes())
	}
}
