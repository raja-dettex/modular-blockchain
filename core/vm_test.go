package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVMRun(t *testing.T) {
	// data := []byte{0x02, 0x0a, 0x061, 0x0c, 0x064, 0x0c, 0x0d}
	data := []byte{0x02, 0x0a, 0x01, 0x0a, 0x0e}
	vm := NewVM(data)
	err := vm.Run()
	assert.Nil(t, err)
	result := vm.stack.Pop().(int)
	assert.Equal(t, result, 1)
	//fmt.Println(string(result))
	// fmt.Printf("stack %v", vm.stack.data)

}

func TestStack(t *testing.T) {
	s := NewStack(128)
	s.Push(1)
	s.Push(2)
	fmt.Println(s.data)
	assert.Equal(t, s.Pop(), 1)
	assert.Equal(t, s.Pop(), 2)
}
