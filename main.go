//usr/bin/go run $0 $@ ; exit
package main

import (
	"bytes"
	"debug/elf"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const versionString = "cdetect 0.3"

var (
	gccMarker   = []byte("GCC: (")
	clangMarker = []byte("clang version")
)

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
		// If the bytes are on this form: "GCC: (GNU) 6.3.0GCC: (GNU) 7.2.0"
		// Then use the last one.
		if bytes.Count(versionData, gccMarker) > 1 {
			// Remove all but the last "GCC: (" version string
			versionData = versionData[bytes.LastIndex(versionData, gccMarker):]
		}
		// Try the first regexp for picking out the version
		versionCatcher1 := regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)\ `)
		gccVersion := bytes.TrimSpace(versionCatcher1.Find(versionData))
		if len(gccVersion) > 0 {
			return "GCC " + string(gccVersion)
		}
		// Try the second regexp for picking out the version
		versionCatcher2 := regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)`)
		gccVersion = bytes.TrimSpace(versionCatcher2.Find(versionData))
		if len(gccVersion) > 0 {
			return "GCC " + string(gccVersion)
		}
		return "GCC " + string(gccVersion)[5:]
	}
	return string(versionData)
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

// returns the compiler name and version that was used for compiling the ELF,
// or an empty string
func compiler(f *elf.File) string {
	if goVersion := gover(f); goVersion != "" {
		return goVersion
	} else if ocamlVersion := ocamlver(f); ocamlVersion != "" {
		return ocamlVersion
	} else if gccVersion := gccver(f); gccVersion != "" {
		return gccVersion
	} else if pasVersion := pasver(f); pasVersion != "" {
		return pasVersion
	} else if tccVersion := tccver(f); tccVersion != "" {
		return tccVersion
	}
	return "unknown"
}

func examine(filename string) string {
	f, err := elf.Open(filename)
	if err != nil {
		if strings.HasPrefix(err.Error(), "bad magic") {
			fmt.Println("Not an ELF")
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	defer f.Close()
	return compiler(f)
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
			fmt.Println(examine(os.Args[1]))
		}
	} else {
		usage()
	}
}
