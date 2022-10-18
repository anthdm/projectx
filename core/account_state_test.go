package core

import (
	"testing"

	"github.com/anthdm/projectx/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAccounState(t *testing.T) {
	state := NewAccountState()

	address := crypto.GeneratePrivateKey().PublicKey().Address()
	account := state.CreateAccount(address)

	assert.Equal(t, account.Address, address)
	assert.Equal(t, account.Balance, uint64(0))

	fetchedAccount, err := state.GetAccount(address)
	assert.Nil(t, err)
	assert.Equal(t, fetchedAccount, account)
}

func TestTransferFailInsufficientBalance(t *testing.T) {
	state := NewAccountState()

	addressBob := crypto.GeneratePrivateKey().PublicKey().Address()
	addressAlice := crypto.GeneratePrivateKey().PublicKey().Address()

	accountBob := state.CreateAccount(addressBob)
	accountBob.Balance = 99

	accountAlice := state.CreateAccount(addressAlice)

	amount := uint64(100)
	assert.NotNil(t, state.Transfer(addressBob, addressAlice, amount))
	assert.Equal(t, accountAlice.Balance, uint64(0))
}

func TestTransferSuccessEmpyToAccount(t *testing.T) {
	state := NewAccountState()

	addressBob := crypto.GeneratePrivateKey().PublicKey().Address()
	addressAlice := crypto.GeneratePrivateKey().PublicKey().Address()

	accountBob := state.CreateAccount(addressBob)
	accountBob.Balance = 100

	accountAlice := state.CreateAccount(addressAlice)

	amount := uint64(100)
	assert.Nil(t, state.Transfer(addressBob, addressAlice, amount))
	assert.Equal(t, accountAlice.Balance, amount)
}
