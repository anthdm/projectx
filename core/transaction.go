package core

import "io"

type Transaction struct {
	Data []byte
}

func (tx *Transaction) DecodeBinary(r io.Reader) error { return nil }

func (tx *Transaction) EncodeBinary(w io.Writer) error { return nil }
