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
var flagsForce bool

func init() {
	flag.BoolVar(&flagsVersion, "V", false, "show version")
	flag.BoolVar(&flagsHelp, "h", false, "show help")
	flag.BoolVar(&flagsForce, "f", false, "delete before format")
	flag.Parse()
}

func main() {
	if flagsVersion {
		fmt.Printf("%s\n", VERSION)
		os.Exit(0)
	}

	if len(flag.Args()) == 1 {
		if err := fs.Make(flag.Arg(0), flagsForce); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}
	} else {
		printHelp()
	}
}

func printHelp() {
	fmt.Printf("%s: [flags] path\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}
