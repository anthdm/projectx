package core

type Storage interface {
	Put(*Block) error
}

type MemoryStore struct {
}

func NewMemorystore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Put(b *Block) error {
	return nil
}
