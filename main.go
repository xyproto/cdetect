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

const versionString = "cdetect 0.5.2"

func usage() {
	fmt.Println(versionString)
	fmt.Println()
	fmt.Println("Detect the compiler version, given an executable (ELF)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    cdetect [OPTION]... [FILE]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("    -v, --version           - version info")
	fmt.Println("    -h, --help              - this help output")
	fmt.Println()
}

// Check if the given filename exists.
// If it exists in $PATH, return the full path,
// else return an empty string.
func which(filename string) (string, error) {
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		return filename, nil
	}
	for _, directory := range strings.Split(os.Getenv("PATH"), ":") {
		fullPath := path.Join(directory, filename)
		_, err := os.Stat(fullPath)
		if !os.IsNotExist(err) {
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
