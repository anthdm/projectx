package core

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/anthdm/projectx/crypto"
	"github.com/stretchr/testify/assert"
)

func TestNFTTransaction(t *testing.T) {
	collectionTx := CollectionTx{
		Fee:      200,
		MetaData: []byte("The beginning of a new collection"),
	}

	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		TxInner: collectionTx,
	}
	tx.Sign(privKey)

	buf := new(bytes.Buffer)
	assert.Nil(t, gob.NewEncoder(buf).Encode(tx))

	txDecoded := &Transaction{}
	assert.Nil(t, gob.NewDecoder(buf).Decode(txDecoded))
	assert.Equal(t, tx, txDecoded)
}

func TestNativeTransferTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		To:    toPrivKey.PublicKey(),
		Value: 666,
	}

	assert.Nil(t, tx.Sign(fromPrivKey))
}

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}

func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(privKey))

	return &tx
}
