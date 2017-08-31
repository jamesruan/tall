package main

import (
	"flag"
	"fmt"
	"github.com/jamesruan/tall/backend/fs"
	"os"
	"path/filepath"
)

const VERSION = "0.0.1"

var flagsVersion bool
var flagsHelp bool

func init() {
	flag.BoolVar(&flagsVersion, "V", false, "show version")
	flag.BoolVar(&flagsHelp, "h", false, "show help")
	flag.Parse()
}

func main() {
	if flagsVersion {
		fmt.Printf("%s\n", VERSION)
		os.Exit(0)
	}

	if len(flag.Args()) == 1 {
		if err := mkfs(flag.Arg(0)); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}
	} else {
		printHelp()
	}
}

func mkfs(entry string) (err error) {
	paths := []string{
		"",
		fs.JOURNALPATH,
		fs.TEMPPATH,
		fs.DATAPATH,
		fs.STATSPATH,
	}
	for _, suffix := range paths {
		path := filepath.Join(entry, suffix)
		fmt.Printf("creating %s\n", path)
		err = os.MkdirAll(path, fs.DefaultDirMode)
		if err != nil {
			return
		}
	}

	path := filepath.Join(entry, fs.SUPERMATPATH)
	fmt.Printf("creating %s\n", path)
	var file *os.File
	file, err = os.Create(path)
	defer file.Close()
	sm := &fs.SuperMeta{
		Level: 1,
	}
	if err = sm.Store(file); err != nil {
		os.RemoveAll(entry)
		return
	}
	return
}

func printHelp() {
	fmt.Printf("%s: path\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}