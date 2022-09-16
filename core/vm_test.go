package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM(t *testing.T) {
	// 1 + 2 = 3
	// 1
	// push stack
	// 2
	// push stack
	// add
	// 3
	// push stack

	data := []byte{0x02, 0x0a, 0x02, 0x0a, 0x0b}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())

	assert.Equal(t, byte(4), vm.stack[vm.sp])
}
