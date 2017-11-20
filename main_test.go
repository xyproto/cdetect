// Test
package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestVoidLinuxNano(t *testing.T) {
	assert.Equal(t, examine("testdata/nano"), "GCC 7.2.0")
}
