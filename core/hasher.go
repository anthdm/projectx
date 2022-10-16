package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/anthdm/projectx/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}

type TxHasher struct{}

// Hash will hash the whole bytes of the TX no exception.
func (TxHasher) Hash(tx *Transaction) types.Hash {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(tx); err != nil {
		panic(err)
	}

	return types.Hash(sha256.Sum256(buf.Bytes()))
}
