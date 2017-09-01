package fs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func StoreStructToFile(s interface{}, f *os.File) (err error) {
	var buf []byte
	buf, err = json.Marshal(s)
	if _, err = f.Seek(0, os.SEEK_SET); err != nil {
		return
	}
	var a int = 0
	var n int = 0

	for a == 0 && a <= len(buf) {
		n, err = f.Write(buf[a:])
		if err != nil {
			return
		}
		a += n
	}
	if err = f.Truncate(int64(len(buf))); err != nil {
		return
	}
	return f.Sync()
}

func LoadStructFromFile(s interface{}, f *os.File) (err error) {
	if _, err = f.Seek(0, os.SEEK_SET); err != nil {
		return
	}
	defer f.Seek(0, os.SEEK_SET)

	var buf []byte
	buf, err = ioutil.ReadAll(f)
	return json.Unmarshal(buf, s)
}
