//usr/bin/go run $0 $@ ; exit
package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/xyproto/ainur"
)

const versionString = "cdetect 0.6.0"

func usage() {
	fmt.Println(versionString + `
Detect the compiler version, given an executable (ELF)

Usage:
    cdetect [OPTION]... [FILE]

Options:
    -v, --version           - version info
    -h, --help              - this help output
	`)
}

// Check if the given filename exists.
// If it exists in $PATH, return the full path,
// else return an empty string.
func which(filename string) (string, error) {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename, nil
	}
	for _, directory := range strings.Split(os.Getenv("PATH"), ":") {
		fullPath := path.Join(directory, filename)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			return fullPath, nil
		}
	}
	return "", errors.New(filename + ": no such file or directory")
}

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-h", "--help":
			usage()
		case "-v", "--version":
			fmt.Println(versionString)
		default:
			filepath, err := which(os.Args[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			compilerVersion, err := ainur.Examine(filepath)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println(compilerVersion)
		}
	} else {
		usage()
	}
}
