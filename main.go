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

const versionString = "cdetect 0.2"

// returns the GCC compiler version or an empty string
// example output: "GCC 6.3.1"
func gccver(f *elf.File) string {
	sec := f.Section(".comment")
	if sec == nil {
		return ""
	}
	b, errData := sec.Data()
	if errData != nil {
		return ""
	}
	cVersion := string(b)
	if strings.Contains(cVersion, "GCC: (") {
		versionCatcher1 := regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)\ `)
		gccVersion := strings.TrimSpace(string(versionCatcher1.Find(b)))
		if gccVersion != "" {
			return "GCC " + gccVersion
		}
		versionCatcher2 := regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)`)
		gccVersion = strings.TrimSpace(string(versionCatcher2.Find(b)))
		if gccVersion != "" {
			return "GCC " + gccVersion
		}
		return "GCC " + cVersion[5:]
	}
	return cVersion
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
	}
	return "unknown"
}

func examine(filename string) {
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
	fmt.Printf("%v\n", compiler(f))
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
}

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			usage()
		} else if os.Args[1] == "-v" || os.Args[1] == "--version" {
			fmt.Println(versionString)
		} else {
			examine(os.Args[1])
		}
	} else {
		usage()
	}
}
