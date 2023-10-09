package types

import (
	"fmt"
	"reflect"
)

type List[T any] struct {
	Data []T
}

func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

func (l *List[T]) Get(index int) T {
	if index > len(l.Data)-1 {
		err := fmt.Errorf("the given index %v is higher than the length %v", index, len(l.Data)-1)
		panic(err)
	}
	return l.Data[index]
}

func (l *List[T]) Insert(v T) {
	l.Data = append(l.Data, v)
}

func (l *List[T]) Clear() {
	l.Data = []T{}
}

func (l *List[T]) GetIndex(v T) int {
	for i := 0; i < len(l.Data); i++ {
		if reflect.DeepEqual(v, l.Data[i]) {
			return i
		}
	}
	return -1
}

func (l *List[T]) Remove(v T) {
	index := l.GetIndex(v)
	if index == -1 {
		return
	}
	l.Pop(index)
}

func (l *List[T]) Pop(index int) {
	l.Data = append(l.Data[:index], l.Data[index+1:]...)
}

func (l *List[T]) Contains(v T) bool {
	for i := 0; i < len(l.Data); i++ {
		if reflect.DeepEqual(v, l.Data[i]) {
			return true
		}
	}
	return false
}

func (l *List[T]) Last() T {
	return l.Data[l.Len()-1]
}

func (l *List[T]) Len() int {
	return len(l.Data)
}
