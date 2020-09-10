package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/takurooo/gopherdojo-studyroom/kadai1/takurooo/transcoder"
)

type Arg struct {
	dir       string
	inFormat  string
	outFormat string
}

var arg Arg

const (
	ExitOK    = 0
	ExitError = 1
)

func init() {
	flag.StringVar(&arg.dir, "d", "", "directory path")
	flag.StringVar(&arg.inFormat, "i", "jpg", "input format(jpg or png)")
	flag.StringVar(&arg.outFormat, "o", "png", "output formatjpg or png")
	flag.Parse()
}

func main() {
	os.Exit(imgtransMain())
}

func imgtransMain() int {
	var err error

	if arg.dir == "" {
		fmt.Fprintln(os.Stderr, "directory path empty")
		flag.Usage()
		return ExitError
	}
	{
		fileInfo, err := os.Stat(arg.dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitError
		}

		if !fileInfo.IsDir() {
			fmt.Fprintf(os.Stderr, "not directory : %s\n", arg.dir)
			return ExitError
		}
	}
	if !transcoder.IsSupported(arg.inFormat) {
		fmt.Fprintf(os.Stderr, "invalid in format : %s\n", arg.inFormat)
		return ExitError
	}
	if !transcoder.IsSupported(arg.outFormat) {
		fmt.Fprintf(os.Stderr, "invalid out format : %s\n", arg.outFormat)
		return ExitError
	}

	trans := transcoder.NewTranscoder(arg.inFormat, arg.outFormat)

	err = filepath.Walk(arg.dir, func(path string, fileInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if !strings.Contains(ext, arg.inFormat) {
			return nil
		}

		err = trans.Do(path)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitError
	}

	return ExitOK
}
