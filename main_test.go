package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestVoidLinuxNano(t *testing.T) {
	assert.Equal(t, mustExamine("testdata/nano_voidlinux"), "GCC 7.2.0")
}

func TestArchLinuxLs(t *testing.T) {
	assert.Equal(t, mustExamine("testdata/ls_archlinux"), "GCC 7.1.1")
}

func TestClang(t *testing.T) {
	assert.Equal(t, mustExamine("testdata/clang_hello"), "Clang 6.0.0")
}

func TestTCC(t *testing.T) {
	assert.Equal(t, mustExamine("testdata/tcc_hello"), "TCC")
}

func TestRustStripped(t *testing.T) {
	assert.Equal(t, mustExamine("testdata/rust_hello_stripped"), "Rust (GCC 8.1.0)")
}

func TestRust(t *testing.T) {
	assert.Equal(t, mustExamine("testdata/rust_hello"), "Rust 1.27.0-nightly")
}

func TestGCC(t *testing.T) {
	assert.Equal(t, mustExamine("testdata/afl-analyze"), "GCC 7.2.0")
}

func TestVersionCompare(t *testing.T) {
	assert.Equal(t, firstIsGreater("2", "1.0.7.abc"), true)
	assert.Equal(t, firstIsGreater("2.0", "2.0 alpha1"), true)
	assert.Equal(t, firstIsGreater("1.0", "2.0.3"), false)
	assert.Equal(t, firstIsGreater("2.0", "2.0.rc1"), true)
}
