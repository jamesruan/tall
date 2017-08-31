// Package fs provides a tall backend based on filesystem
package fs

import (
	"bytes"
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

// Store save file to a tempfile, calculate score for the data
// and move to the file with the score as its name.
// The path for the file is scattered according to SuperMeta.Level
func (b *FSBackend) Store(data []byte) (score tall.HexBytes, err error) {
	var agent *FSWriteAgent
	r := bytes.NewBuffer(data)
	if agent, err = b.NewWriteAgent(); err != nil {
		return
	}

	if _, err = io.Copy(agent, r); err != nil {
		return
	}
	if err = agent.Close(); err != nil {
		return
	}
	score = agent.Score()

	return score, nil
}

func (b *FSBackend) NewWriteAgent() (*FSWriteAgent, error) {
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

func (b *FSBackend) Load(score tall.HexBytes) (data []byte, err error) {
	buf := new(bytes.Buffer)
	var r io.ReadCloser
	r, err = b.Open(score)
	defer r.Close()
	if _, err = io.Copy(buf, r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *FSBackend) Open(score tall.HexBytes) (*os.File, error) {
	var file *os.File
	var err error
	path := filepath.Join(b.entry, scoreToPath(score, b.sm.Level))
	if file, err = os.Open(path); err != nil {
		return nil, err
	}
	return file, nil
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
