package core

import (
	"encoding/binary"
	"fmt"
)

type Instruction byte

const (
	InstrPushInt  Instruction = 0x0a
	InstrAdd      Instruction = 0x0b
	InstrPushByte Instruction = 0x0c
	InstrPack     Instruction = 0x0d
	InstrSub      Instruction = 0x0e
	InstrStore    Instruction = 0x0f
	InstrGet      Instruction = 0xae
	InstrMul      Instruction = 0xea
	InstrDiv      Instruction = 0xdf
)

type Stack struct {
	data []any
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
	}
}
func (s *Stack) Push(value any) {
	s.data = append([]any{value}, s.data...)
}
func (s *Stack) Pop() any {
	val := s.data[0]
	s.data = s.data[1:]
	return val
}

type VM struct {
	data          []byte
	ip            int
	stack         *Stack
	ContractState *State
}

func NewVM(data []byte, contractState *State) *VM {
	return &VM{
		data:          data,
		ip:            0,
		stack:         NewStack(128),
		ContractState: contractState,
	}
}

func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])
		err := vm.Exec(instr)
		if err != nil {
			return err
		}
		vm.ip++
		if vm.ip > len(vm.data)-1 {
			break
		}
	}
	return nil
}

func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrGet:
		key := vm.stack.Pop().([]byte)
		b, err := vm.ContractState.Get(key)
		if err != nil {
			return err
		}
		// buff := make([]byte, len(b))
		// for _, val := range b {
		// 	buff = append(buff, val)
		// }
		fmt.Printf("value Retrieved : %v\n", b)
		vm.stack.Push(b)

	case InstrStore:
		var serializedValue []byte
		key := vm.stack.Pop().([]byte)
		val := vm.stack.Pop()
		switch v := val.(type) {
		case int:
			serializedValue = serializeInt64(int64(v))
			vm.ContractState.Put(key, serializedValue)
		default:
			panic("Default TODO //")
		}

	case InstrPushInt:
		vm.stack.Push(int(vm.data[vm.ip-1]))
	case InstrAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a + b
		vm.stack.Push(c)
	case InstrPushByte:
		vm.stack.Push(byte(vm.data[vm.ip-1]))
	case InstrPack:
		n := vm.stack.Pop().(int)
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = vm.stack.Pop().(byte)
		}
		vm.stack.Push(b)
	case InstrSub:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a - b
		vm.stack.Push(c)
	case InstrMul:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a * b
		vm.stack.Push(c)
	case InstrDiv:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a / b
		vm.stack.Push(c)
	}
	return nil
}

// func (vm *VM) firstNode() any {
// 	return vm.stack.data[vm.stack.sp]
// }

func serializeInt64(value int64) []byte {
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, uint64(value))
	return buff
}

func deserializeInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}
