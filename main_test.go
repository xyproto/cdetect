// Test
package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestVoidLinuxNano(t *testing.T) {
	assert.Equal(t, examine("testdata/nano_voidlinux"), "GCC 7.2.0")
}

func TestArchLinuxLs(t *testing.T) {
	assert.Equal(t, examine("testdata/ls_archlinux"), "GCC 7.1.1")
}

func TestClang(t *testing.T) {
	assert.Equal(t, examine("testdata/clang_hello"), "Clang 5.0.0")
}

func TestTCC(t *testing.T) {
	assert.Equal(t, examine("testdata/tcc_hello"), "TCC")
}
