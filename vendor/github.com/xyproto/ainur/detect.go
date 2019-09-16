// Package ainur provides functions for examining ELF files
package ainur

// Only depends on the standard library
import (
	"bytes"
	"debug/elf"
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	gccMarker   = []byte("GCC: (")
	gnuEnding   = []byte("GNU) ")
	clangMarker = []byte("clang version")
	rustMarker  = []byte("rustc version")
	ghcMarker   = []byte("GHC ")
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

// FirstIsGreater checks if the first version number is greater than the second one.
// It uses a relatively simple algorithm, where all non-numbers counts as less than "0".
func FirstIsGreater(a, b string) bool {
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

// GHCVer returns the GHC compiler version or an empty string
// example output: "GHC 8.6.2"
func GHCVer(f *elf.File) string {
	sec := f.Section(".comment")
	if sec == nil {
		return ""
	}
	versionData, errData := sec.Data()
	if errData != nil {
		return ""
	}
	if bytes.Contains(versionData, ghcMarker) {
		// Try the first regexp for picking out the version
		versionCatcher1 := regexp.MustCompile(`GHC\ (\d{1,4}\.)(\d+\.)?(\d+)`)
		ghcVersion := bytes.TrimSpace(versionCatcher1.Find(versionData))
		if len(ghcVersion) > 0 {
			return "GHC " + string(ghcVersion[4:])
		}
	}
	return ""
}

// GCCVer returns the GCC compiler version or an empty string
// example output: "GCC 6.3.1"
// Also handles clang.
func GCCVer(f *elf.File) string {
	debug := false
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
			if FirstIsGreater(string(versionA), string(versionB)) {
				versionData = versionA
			} else {
				versionData = versionB
			}
		}
		// Try the first regexp for picking out the version
		versionCatcher1 := regexp.MustCompile(`\) (\d{1,4}\.)(\d+\.)?(\*|\d+)\ `)
		gccVersion := bytes.TrimSpace(versionCatcher1.Find(versionData))
		if len(gccVersion) > 0 {
			if debug {
				println("GCC #1 " + string(gccVersion[2:]))
			}
			return "GCC " + string(gccVersion[2:])
		}
		// Try the second regexp for picking out the version
		versionCatcher2 := regexp.MustCompile(` (\d{1,4}\.)(\d+\.)?(\*|\d+)`)
		gccVersion = bytes.TrimSpace(versionCatcher2.Find(versionData))
		if len(gccVersion) > 0 {
			if debug {
				println("GCC #2 " + string(gccVersion))
			}
			// Check that it does not start with "1.", that may happen
			if !bytes.HasPrefix(gccVersion, []byte("1.")) {
				return "GCC " + string(gccVersion)
			}
		}
		// Try the third regexp for picking out the version
		versionCatcher3 := regexp.MustCompile(`(\d{1,4}\.)(\d+\.)?(\*|\d+)`)
		gccVersion = bytes.TrimSpace(versionCatcher3.Find(versionData))
		if len(gccVersion) > 0 {
			if debug {
				println("GCC #3 " + string(gccVersion))
			}
			// Check that it does not start with "1.", that may happen
			if !bytes.HasPrefix(gccVersion, []byte("1.")) {
				return "GCC " + string(gccVersion)
			}
		}
		// Try the fourth regexp for picking out the version
		versionCatcher4 := regexp.MustCompile(`\) (\d{1,4}\.)(\d+\.)?(\*|\d+).(\d+)`)
		gccVersion = bytes.TrimSpace(versionCatcher4.Find(versionData))
		if len(gccVersion) > 0 {
			if debug {
				println("GCC #4 " + string(gccVersion))
			}
			return "GCC " + string(gccVersion)[2:]
		}
		// See what we've got
		gccVersionString := strings.TrimSpace(string(gccVersion))
		if len(gccVersionString) > 5 {
			if debug {
				println("GCC #4 " + string(gccVersion[5:]))
			}
			// Check that the version number is not "0"
			retver := string(gccVersion)[5:]
			if retver != "0" {
				return "GCC " + retver
			}
		}
		// Failed to find a GCC version string
		return ""
	}
	return string(versionData)
}

// RustVerUnstripped returns the Rust compiler version or an empty string
// example output: "Rust 1.27.0"
func RustVerUnstripped(f *elf.File) string {
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

// RustVerStripped returns the Rust compiler version or an empty string,
// from a stripped Rust executable. Does not contain the Rust
// version number.
// Example output: "Rust (GCC 8.1.0)"
func RustVerStripped(f *elf.File) string {
	// Check if the .gcc_except_table ELF section exists
	if f.Section(".gcc_except_table") == nil {
		return ""
	}
	// Check if the .rodata ELF section exists
	sec := f.Section(".rodata")
	if sec == nil {
		return ""
	}
	b, errData := sec.Data()
	if errData != nil {
		return ""
	}
	// Look for the rust marker that may appear in new, stripped executables
	if bytes.Index(b, []byte("/rustc-")) == -1 {
		// Look for the rust marker that may appear in old, stripped executables
		rustIndex1 := bytes.Index(b, []byte("__rust_"))
		if rustIndex1 <= 0 || b[rustIndex1-1] != 0 {
			// No rust markers! Probably not created with the Rust compiler.
			return ""
		}
	}
	// Rust may use GCC for linking
	if gccVersion := GCCVer(f); gccVersion != "" {
		return "Rust (" + GCCVer(f) + ")"
	}
	return "Rust"
}

// DVer returns "DMD" if it is detected
// Example output: "DMD"
func DVer(f *elf.File) string {
	// Check if the .dynstr ELF section exists
	sec := f.Section(".dynstr")
	if sec == nil {
		return ""
	}
	b, errData := sec.Data()
	if errData != nil {
		return ""
	}
	// Look for the DMD marker
	if bytes.Index(b, []byte("__dmd_")) != -1 {
		return "DMD"
	}
	return ""
}

// GoVer returns the Go compiler version or an empty string
// example output: "Go 1.8.3"
func GoVer(f *elf.File) string {
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

// PasVer returns the FPC compiler version or an empty string
// example output: "FPC 3.0.2"
func PasVer(f *elf.File) string {
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

// TCCVer returns "TCC" or an empty string
// TCC has no version number, but it does have some signature sections.
func TCCVer(f *elf.File) string {
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

// OCamlVer returns the OCaml compiler version or an empty string
// example output: "OCaml 4.05.0"
func OCamlVer(f *elf.File) string {
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

// Compiler takes an *elf.File and tries to find which compiler and version
// it was compiled with, by probing for known locations, strings and patterns.
func Compiler(f *elf.File) string {
	// The ordering matters
	if goVersion := GoVer(f); goVersion != "" {
		return goVersion
	} else if ocamlVersion := OCamlVer(f); ocamlVersion != "" {
		return ocamlVersion
	} else if ghcVersion := GHCVer(f); ghcVersion != "" {
		return ghcVersion
	} else if rustVersion := RustVerUnstripped(f); rustVersion != "" {
		return rustVersion
	} else if rustVersion := RustVerStripped(f); rustVersion != "" {
		return rustVersion
	} else if dVersion := DVer(f); dVersion != "" {
		return dVersion
	} else if gccVersion := GCCVer(f); gccVersion != "" {
		return gccVersion
	} else if pasVersion := PasVer(f); pasVersion != "" {
		return pasVersion
	} else if tccVersion := TCCVer(f); tccVersion != "" {
		return tccVersion
	}
	return "unknown"
}

// Stripped returns true if symbols can not be retrieved from the given ELF file
func Stripped(f *elf.File) bool {
	_, err := f.Symbols()
	return err != nil
}

// Examine tries to discover which compiler and compiler version the given
// file was compiled with.
func Examine(filename string) (string, error) {
	f, err := elf.Open(filename)
	if err != nil {
		if strings.HasPrefix(err.Error(), "bad magic") {
			return "", errors.New(filename + ": Not an ELF")
		}
		return "", err
	}
	defer f.Close()
	return Compiler(f), nil
}

// MustExamine does the same as examine, but panics instead of returning an error
func MustExamine(filename string) string {
	compilerVersion, err := Examine(filename)
	if err != nil {
		panic(err)
	}
	return compilerVersion
}
