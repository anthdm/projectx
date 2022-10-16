package core

import (
	"fmt"
	"sync"

	"github.com/anthdm/projectx/types"
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

func (s *AccountState) Transfer(from, to types.Address, amount uint64) error {
	if err := s.SubBalance(from, amount); err != nil {
		return err
	}

	return s.AddBalance(to, amount)
}

func (s *AccountState) GetBalance(to types.Address) (uint64, error) {
	balance, ok := s.state[to]
	if !ok {
		return 0.0, fmt.Errorf("address (%s) unkown", to)
	}

	return balance, nil
}

func (s *AccountState) SubBalance(to types.Address, amount uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	balance, ok := s.state[to]
	if !ok {
		return fmt.Errorf("address (%s) unknown", to)
	}

	if balance < amount {
		return fmt.Errorf("insuccient account balance (%d) for amount (%d)", balance, amount)
	}

	s.state[to] -= amount

	return nil
}

func (s *AccountState) AddBalance(to types.Address, amount uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// _, ok := s.state[to]
	s.state[to] += amount
	// if !ok {
	// 	s.state[to] = amount
	// } else {
	// 	s.state[to] += amount
	// }

	return nil
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

func (s *State) Get(k []byte) ([]byte, error) {
	key := string(k)

	value, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("given key %s not found", key)
	}

	return value, nil
}
