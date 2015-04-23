package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

// Everything that we can't determine as a directory
// is not one (files, links, ...)
func isRealDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return false
	}
	if s.IsDir() && (s.Mode()&os.ModeSymlink == 0) {
		return true
	}
	return false
}

func myRmDir(dir string) error {
	err := syscall.Rmdir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't rmdir %s reason: %s\n", dir, err)
	}
	return err
}

func rmdirR(directory string) {
	if isRealDir(directory) {
		files, _ := ioutil.ReadDir(directory)
		for _, basename := range files {
			file := directory + "/" + basename.Name()
			rmdirR(file)
			//fmt.Printf("Debug: should rmdir %s\n", file)
		}
		myRmDir(directory)
	} else {
		fmt.Fprintf(os.Stderr, "file: %s is not a directory\n", directory)
		os.Exit(2)
	}
}

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		fmt.Printf("rmdir_recurs. %s\n", arg)
		rmdirR(arg)
	}
}

// vim: set ts=2 sw=2 list ft=go:
