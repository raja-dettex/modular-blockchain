package core

import (
	"fmt"
	"sync"

	"github.com/raja-dettex/modular-blockchain/types"
)

type AccountState struct {
	mu    sync.RWMutex
	state map[types.Address]uint64
}

func NewAccountState() *AccountState {
	return &AccountState{
		state: make(map[types.Address]uint64),
	}
}

func (as *AccountState) AddBalance(to types.Address, amount uint64) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	_, ok := as.state[to]
	if !ok {
		as.state[to] = amount
	} else {
		as.state[to] += amount
	}
	//fmt.Printf("money added to account of addresss %v\n", to)
	return nil
}
func (as *AccountState) DeductBalance(from types.Address, amount uint64) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	balance, ok := as.state[from]
	fmt.Println(ok)
	if !ok {
		return fmt.Errorf("account state ( %v ) does not exist", from)
	}
	if balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	as.state[from] -= amount

	return nil
}

func (as *AccountState) GetBalance(from types.Address) (uint64, error) {
	balance, ok := as.state[from]
	if !ok {
		return 0.0, fmt.Errorf("accont state %s does not exist", from)
	}
	return balance, nil
}

func (as *AccountState) Transfer(from, to types.Address, amount uint64) error {
	if err := as.DeductBalance(from, amount); err != nil {
		return err
	}
	return as.AddBalance(to, amount)
}

type State struct {
	data map[string][]byte
}

func NewState() *State {
	return &State{
		data: make(map[string][]byte),
	}
}

func (s *State) Put(k, v []byte) error {
	s.data[string(k)] = v
	return nil
}

func (s *State) Delete(k []byte) error {
	delete(s.data, string(k))
	return nil
}

func (s *State) Get(key []byte) ([]byte, error) {
	val, ok := s.data[string(key)]
	if !ok {
		return nil, fmt.Errorf("key: %v is not in the state ", string(key))
	}
	return val, nil
}
