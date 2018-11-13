//usr/bin/go run $0 $@ ; exit
package main

// Only depends on the standard library
import (
	"bytes"
	"debug/elf"
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const versionString = "cdetect 0.5.1"

var (
	gccMarker   = []byte("GCC: (")
	gnuEnding   = []byte("GNU) ")
	clangMarker = []byte("clang version")
	rustMarker  = []byte("rustc version")
)

// versionSum takes a slice of strings that are the parts of a version number.
// The parts are converted to numbers. If they can't be converted, they count
// as less than nothing. The parts are then summed together, but with more
// emphasis put on the earlier numbers. 2.0.0.0 has emphasis 2000.
// The sum is then returned.
func versionSum(parts []string) int {
	sum := 0
	length := len(parts)
	for i := length - 1; i >= 0; i-- {
		num, err := strconv.Atoi(parts[i])
		if err != nil {
			num = -1
		}
		sum += num * int(math.Pow(float64(10), float64(length-i-1)))
	}
	return sum
}

// firstIsGreater checks if the first version number is greater than the second one.
// It uses a relatively simple algorithm, where all non-numbers counts as less than "0".
func firstIsGreater(a, b string) bool {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	// Expand the shortest version list with zeroes
	for len(aParts) < len(bParts) {
		aParts = append(aParts, "0")
	}
	for len(bParts) < len(aParts) {
		bParts = append(bParts, "0")
	}
	// The two lists that are being compared should be of the same length
	return versionSum(aParts) > versionSum(bParts)
}

// returns the GCC compiler version or an empty string
// example output: "GCC 6.3.1"
// Also handles clang.
func gccver(f *elf.File) string {
	sec := f.Section(".comment")
	if sec == nil {
		return ""
	}
	versionData, errData := sec.Data()
	if errData != nil {
		return ""
	}
	if bytes.Contains(versionData, gccMarker) {
		// Check if this is really clang
		if bytes.Contains(versionData, clangMarker) {
			clangVersionCatcher := regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)\ `)
			clangVersion := bytes.TrimSpace(clangVersionCatcher.Find(versionData))
			return "Clang " + string(clangVersion)
		}
		// If the bytes are on this form: "GCC: (GNU) 6.3.0GCC: (GNU) 7.2.0",
		// use the largest version number.
		if bytes.Count(versionData, gccMarker) > 1 {
			// Split in to 3 parts, always valid for >=2 instances of gccMarker
			elements := bytes.SplitN(versionData, gccMarker, 3)
			versionA := elements[1]
			versionB := elements[2]
			if bytes.HasPrefix(versionA, gnuEnding) {
				versionA = versionA[5:]
			}
			if bytes.HasPrefix(versionB, gnuEnding) {
				versionB = versionB[5:]
			}
			if firstIsGreater(string(versionA), string(versionB)) {
				versionData = versionA
			} else {
				versionData = versionB
			}
		}
		// Try the first regexp for picking out the version
		versionCatcher1 := regexp.MustCompile(`(\d{1,4}\.)(\d+\.)?(\*|\d+)\ `)
		gccVersion := bytes.TrimSpace(versionCatcher1.Find(versionData))
		if len(gccVersion) > 0 {
			return "GCC " + string(gccVersion)
		}
		// Try the second regexp for picking out the version
		versionCatcher2 := regexp.MustCompile(` (\d{1,4}\.)(\d+\.)?(\*|\d+)`)
		gccVersion = bytes.TrimSpace(versionCatcher2.Find(versionData))
		if len(gccVersion) > 0 {
			return "GCC " + string(gccVersion)
		}
		// Try the third regexp for picking out the version
		versionCatcher3 := regexp.MustCompile(`(\d{1,4}\.)(\d+\.)?(\*|\d+)`)
		gccVersion = bytes.TrimSpace(versionCatcher3.Find(versionData))
		if len(gccVersion) > 0 {
			return "GCC " + string(gccVersion)
		}
		// See what we've got
		gccVersionString := strings.TrimSpace(string(gccVersion))
		if len(gccVersionString) > 5 {
			return "GCC " + string(gccVersion)[5:]
		}
		// Failed to find a GCC version string
		return ""
	}
	return string(versionData)
}

// Returns the Rust compiler version or an empty string
// example output: "Rust 1.27.0"
func rustverUnstripped(f *elf.File) string {
	// Check if there is debug data in the executable, that may contain the version number
	sec := f.Section(".debug_str")
	if sec == nil {
		return ""
	}
	b, errData := sec.Data()
	if errData != nil {
		return ""
	}

	pos1 := bytes.Index(b, rustMarker)
	if pos1 == -1 {
		return ""
	}
	pos1 += len(rustMarker) + 1
	pos2 := bytes.Index(b[pos1:], []byte("("))
	if pos2 == -1 {
		return ""
	}
	pos2 += pos1
	versionString := strings.TrimSpace(string(b[pos1:pos2]))

	return "Rust " + versionString
}

// Returns the Rust compiler version or an empty string,
// from a stripped Rust executable. Does not contain the Rust
// version number.
// Example output: "Rust (GCC 8.1.0)"
func rustverStripped(f *elf.File) string {
	sec := f.Section(".rodata")
	if sec == nil {
		return ""
	}
	b, errData := sec.Data()
	if errData != nil {
		return ""
	}
	foundIndex := bytes.Index(b, []byte("__rust_"))
	if foundIndex <= 0 || b[foundIndex-1] != 0 {
		return ""
	}
	// Rust may use GCC for linking
	if gccVersion := gccver(f); gccVersion != "" {
		return "Rust (" + gccver(f) + ")"
	}
	return "Rust"
}

// returns the Go compiler version or an empty string
// example output: "Go 1.8.3"
func gover(f *elf.File) string {
	sec := f.Section(".rodata")
	if sec == nil {
		return ""
	}
	b, errData := sec.Data()
	if errData != nil {
		return ""
	}
	versionCatcher := regexp.MustCompile(`go(\d+\.)(\d+\.)?(\*|\d+)`)
	goVersion := string(versionCatcher.Find(b))
	if strings.HasPrefix(goVersion, "go") {
		return "Go " + goVersion[2:]
	}
	if goVersion == "" {
		gosec := f.Section(".gosymtab")
		if gosec != nil {
			return "Go (unknown version)"
		}
		return ""
	}
	return goVersion
}

// returns the FPC compiler version or an empty string
// example output: "FPC 3.0.2"
func pasver(f *elf.File) string {
	sec := f.Section(".data")
	if sec == nil {
		return ""
	}
	b, errData := sec.Data()
	if errData != nil {
		return ""
	}
	versionCatcher := regexp.MustCompile(`FPC\ (\d+\.)?(\d+\.)?(\*|\d+)`)
	return string(versionCatcher.Find(b))

}

// TCC has no version number, but it has some signature sections
// Returns "TCC" or an empty string
func tccver(f *elf.File) string {
	// .note.ABI-tag must be missing
	if f.Section(".note.ABI-tag") != nil {
		// TCC does not normally have this section, not TCC
		return ""
	}
	if f.Section(".rodata.cst4") == nil {
		// TCC usually has this section, not TCC
		return ""
	}
	return "TCC"
}

// returns the OCaml compiler version or an empty string
// example output: "OCaml 4.05.0"
func ocamlver(f *elf.File) string {
	sec := f.Section(".rodata")
	if sec == nil {
		return ""
	}
	b, errData := sec.Data()
	if errData != nil {
		return ""
	}
	if !bytes.Contains(b, []byte("[ocaml]")) {
		// Probably not OCaml
		return ""
	}
	versionCatcher := regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)`)
	ocamlVersion := "OCaml " + string(versionCatcher.Find(b))
	if ocamlVersion == "" {
		return "OCaml (unknown version)"
	}
	return ocamlVersion
}

// compiler takes an *elf.File and tries to find which compiler and version
// it was compiled with, by probing for known locations, strings and patterns.
func compiler(f *elf.File) string {
	// The ordering matters
	if goVersion := gover(f); goVersion != "" {
		return goVersion
	} else if ocamlVersion := ocamlver(f); ocamlVersion != "" {
		return ocamlVersion
	} else if rustVersion := rustverUnstripped(f); rustVersion != "" {
		return rustVersion
	} else if rustVersion := rustverStripped(f); rustVersion != "" {
		return rustVersion
	} else if gccVersion := gccver(f); gccVersion != "" {
		return gccVersion
	} else if pasVersion := pasver(f); pasVersion != "" {
		return pasVersion
	} else if tccVersion := tccver(f); tccVersion != "" {
		return tccVersion
	}
	return "unknown"
}

// examine tries to discover which compiler and compiler version the given
// file was compiled with.
func examine(filename string) (string, error) {
	f, err := elf.Open(filename)
	if err != nil {
		if strings.HasPrefix(err.Error(), "bad magic") {
			return "", errors.New(filename + ": Not an ELF")
		}
		return "", err
	}
	defer f.Close()
	return compiler(f), nil
}

// mustExamine does the same as examine, but panics instead of returning an error
func mustExamine(filename string) string {
	compilerVersion, err := examine(filename)
	if err != nil {
		panic(err)
	}
	return compilerVersion
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
			compilerVersion, err := examine(filepath)
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
