package fs

import (
	"os"
	"testing"
)

func TestMake(t *testing.T) {
	path := "test"
	if err := Make(path, true); err != nil {
		t.Error(err)
	}
	if err := Make(path, false); err == nil {
		t.Error("exist detect")
	}
	if err := Make(path, true); err != nil {
		t.Error(err)
	}
	os.RemoveAll(path)
}
