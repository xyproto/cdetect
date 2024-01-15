package ainur

import (
	"debug/elf"
	"errors"
)

// Static checks that PT_DYNAMIC is not in one of the program headers of the ELF file
func Static(f *elf.File) bool {
	for _, prog := range f.Progs {
		progType := prog.ProgHeader.Type
		if progType == elf.PT_DYNAMIC {
			return false
		}
	}
	return true
}

// ExamineStatic opens the given filename and checks that it is an ELF file.
// It then calls Static to confirm that PT_DYNAMIC is not present in the program headers.
func ExamineStatic(filename string) (bool, error) {
	f, err := elf.Open(filename)
	if err != nil {
		if _, isFormatError := err.(*elf.FormatError); isFormatError {
			return false, errors.New(filename + ": Not an ELF")
		}
		return false, err
	}
	defer f.Close()

	// This is where the actual
	return Static(f), nil
}

// MustExamineStatic does the same as ExamineStatic, but panics instead of returning an error
func MustExamineStatic(filename string) bool {
	static, err := ExamineStatic(filename)
	if err != nil {
		panic(err)
	}
	return static
}
