// Package fs provides a tall backend based on filesystem
package fs

import (
	"github.com/jamesruan/tall"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// implements tall.Backend, tall.BackendScavenger
type FSBackend struct {
	entry  string
	sm     *SuperMeta
	logger *log.Logger
	debug  bool
	lock   *sync.Mutex
}

func New(entry string) (*FSBackend, error) {
	var supermetafile *os.File
	var err error
	path := filepath.Join(entry, SUPERMATPATH)
	supermetafile, err = os.OpenFile(path, os.O_RDWR, DefaultMode)
	if err != nil {
		return nil, err
	}
	sm := new(SuperMeta)
	if err = sm.Load(supermetafile); err != nil {
		supermetafile.Close()
		return nil, err
	}
	be := &FSBackend{
		entry: entry,
		sm:    sm,
		lock:  new(sync.Mutex),
	}
	return be, nil
}

func (b *FSBackend) Score(d []byte) tall.HexBytes {
	return Score(d)
}

func (b *FSBackend) Create() (*FSWriteAgent, error) {
	var tempfile *os.File
	var err error

	// TODO: use journal
	if tempfile, err = ioutil.TempFile(filepath.Join(b.entry, TEMPPATH), "data"); err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	w := io.MultiWriter(tempfile, pw)

	return &FSWriteAgent{
		b:        b,
		sch:      ScoreFrom(pr),
		tempfile: tempfile,
		w:        w,
		pw:       pw,
		pr:       pr,
	}, nil
}

func (b *FSBackend) Open(score tall.HexBytes) (*FSReadAgent, error) {
	var file *os.File
	var err error
	path := filepath.Join(b.entry, scoreToPath(score, b.sm.Level))
	if file, err = os.Open(path); err != nil {
		return nil, err
	}
	return &FSReadAgent{file}, nil
}

func (b *FSBackend) Has(score tall.HexBytes) (ok bool, err error) {
	path := filepath.Join(b.entry, scoreToPath(score, b.sm.Level))
	if _, err = os.Stat(path); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (b *FSBackend) SetLogger(logger *log.Logger, debug bool) {
	b.lock.Lock()
	b.logger = logger
	b.debug = debug
	b.lock.Unlock()
}

func scoreToPath(score tall.HexBytes, level int) string {
	block := 2
	scoreleft := string(score)
	index := 0
	arr := []string{}
	for index < level {
		var v string
		v, scoreleft = scoreleft[:block*index], scoreleft[block*index:]
		arr = append(arr, v)
		index += 1
	}
	arr = append(arr, string(score))
	return filepath.Join(arr...)
}

func getScoreAndWrite(f io.WriteCloser, r io.Reader) (tall.HexBytes, error) {
	var err error
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()
	writer := io.MultiWriter(f, pw)

	score := ScoreFrom(pr)

	if _, err = io.Copy(writer, r); err != nil {
		return "", err
	}
	return <-score, err
}
