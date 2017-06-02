//usr/bin/go run $0 $@ ; exit
package main

import (
	"debug/elf"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const versionString = "ELFinfo 0.3"

// stripped returns true if symbols can not be retrieved from the given ELF file
func stripped(f *elf.File) bool {
	_, err := f.Symbols()
	return err != nil
}

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
	if strings.Contains(cVersion, "GCC: (GNU) ") {
		versionCatcher := regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)`)
		gccVersion := string(versionCatcher.Find(b))
		if gccVersion != "" {
			return "GCC " + gccVersion
		}
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

// returns the compiler name and version that was used for compiling the ELF,
// or an empty string
func compiler(f *elf.File) string {
	if goVersion := gover(f); goVersion != "" {
		return goVersion
	}
	if gccVersion := gccver(f); gccVersion != "" {
		return gccVersion
	}
	if pasVersion := pasver(f); pasVersion != "" {
		return pasVersion
	}
	return "unknown"
}

func machine2string(m elf.Machine) string {
	switch m {
	case elf.EM_NONE:
		return "Unknown machine"
	case elf.EM_M32:
		return "AT&T WE32100"
	case elf.EM_SPARC:
		return "Sun SPARC"
	case elf.EM_386:
		return "Intel i386"
	case elf.EM_68K:
		return "Motorola 68000"
	case elf.EM_88K:
		return "Motorola 88000"
	case elf.EM_860:
		return "Intel i860"
	case elf.EM_MIPS:
		return "MIPS R3000 Big-Endian only"
	case elf.EM_S370:
		return "IBM System/370"
	case elf.EM_MIPS_RS3_LE:
		return "MIPS R3000 Little-Endian"
	case elf.EM_PARISC:
		return "HP PA-RISC"
	case elf.EM_VPP500:
		return "Fujitsu VPP500"
	case elf.EM_SPARC32PLUS:
		return "SPARC v8plus"
	case elf.EM_960:
		return "Intel 80960"
	case elf.EM_PPC:
		return "PowerPC 32-bit"
	case elf.EM_PPC64:
		return "PowerPC 64-bit"
	case elf.EM_S390:
		return "IBM System/390"
	case elf.EM_V800:
		return "NEC V800"
	case elf.EM_FR20:
		return "Fujitsu FR20"
	case elf.EM_RH32:
		return "TRW RH-32"
	case elf.EM_RCE:
		return "Motorola RCE"
	case elf.EM_ARM:
		return "ARM"
	case elf.EM_SH:
		return "Hitachi SH"
	case elf.EM_SPARCV9:
		return "SPARC v9 64-bit"
	case elf.EM_TRICORE:
		return "Siemens TriCore embedded processor"
	case elf.EM_ARC:
		return "Argonaut RISC Core"
	case elf.EM_H8_300:
		return "Hitachi H8/300"
	case elf.EM_H8_300H:
		return "Hitachi H8/300H"
	case elf.EM_H8S:
		return "Hitachi H8S"
	case elf.EM_H8_500:
		return "Hitachi H8/500"
	case elf.EM_IA_64:
		return "Intel IA-64 Processor"
	case elf.EM_MIPS_X:
		return "Stanford MIPS-X"
	case elf.EM_COLDFIRE:
		return "Motorola ColdFire"
	case elf.EM_68HC12:
		return "Motorola M68HC12"
	case elf.EM_MMA:
		return "Fujitsu MMA"
	case elf.EM_PCP:
		return "Siemens PCP"
	case elf.EM_NCPU:
		return "Sony nCPU"
	case elf.EM_NDR1:
		return "Denso NDR1 microprocessor"
	case elf.EM_STARCORE:
		return "Motorola Star*Core processor"
	case elf.EM_ME16:
		return "Toyota ME16 processor"
	case elf.EM_ST100:
		return "STMicroelectronics ST100 processor"
	case elf.EM_TINYJ:
		return "Advanced Logic Corp. TinyJ processor"
	case elf.EM_X86_64:
		return "Advanced Micro Devices x86-64"
	case elf.EM_AARCH64:
		return "ARM 64-bit Architecture (AArch64)"
	}
	return "Unknown machine"
}

func examine(filename string) {
	f, err := elf.Open(filename)
	if err != nil {
		fmt.Printf("%s: not an ELF: %s\n", filename, err.Error())

		os.Exit(1)
	}
	fmt.Printf("%s: stripped=%v, compiler=%v, byteorder=%v, machine=%v\n", filename, stripped(f), compiler(f), f.ByteOrder, machine2string(f.Machine))
	f.Close()
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" {
			fmt.Println(versionString)
			return
		}
		examine(os.Args[1])
	} else {
		fmt.Println("Needs a filename as the first argument")
		os.Exit(1)
	}
}
