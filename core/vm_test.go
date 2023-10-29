package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVMRun(t *testing.T) {
	// data := []byte{0x02, 0x0a, 0x061, 0x0c, 0x064, 0x0c, 0x0d}

	state := NewState()
	vm := NewVM(Contract(), state)
	err := vm.Run()
	assert.Nil(t, err)
	fmt.Println(vm.stack.data)
	fmt.Println(state)
	//val, err := state.Get([]byte("ad"))
	//assert.Nil(t, err)
	//assert.Equal(t, deserializeInt64(val), int64(5))
	//result := vm.stack.Pop().(int)
	//assert.Equal(t, result, 1)
	// fmt.Println(result)
	// fmt.Printf("stack %v", vm.stack.data)

}

func TestInstrMul(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x04, 0x0a, 0xea}
	state := NewState()
	vm := NewVM(data, state)
	err := vm.Run()
	assert.Nil(t, err)
	result := vm.stack.Pop()
	assert.Equal(t, result, 12)
}
func TestInstrDiv(t *testing.T) {
	data := []byte{0x02, 0x0a, 0x04, 0x0a, 0xdf}
	state := NewState()
	vm := NewVM(data, state)
	err := vm.Run()
	assert.Nil(t, err)
	result := vm.stack.Pop()
	assert.Equal(t, result, 2)
}

func Contract() []byte {
	keyBytes := []byte{0x061, 0x0c, 0x064, 0x0c, 0x02, 0x0a, 0x0d, 0x0ae}
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0e, 0x061, 0x0c, 0x064, 0x0c, 0x02, 0x0a, 0x0d, 0x0f}
	data = append(data, keyBytes...)
	return data
}
func TestStack(t *testing.T) {
	s := NewStack(128)
	s.Push(1)
	s.Push(2)
	fmt.Println(s.data)
	assert.Equal(t, s.Pop(), 1)
	assert.Equal(t, s.Pop(), 2)
}
