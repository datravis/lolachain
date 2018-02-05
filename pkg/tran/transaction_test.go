package tran

import (
	"testing"
	"time"

	"github.com/datravis/lolachain/pkg/keys"

	"github.com/stretchr/testify/assert"
)

// TestNewTransaction verifies we're able to use the provided helper to create a transaction.
func TestNewTransaction(t *testing.T) {
	sym := "TEST"
	source := "source_address"
	dest := "dest_address"
	amount := 1.0
	memo := "memo"
	tm := time.Unix(0, 0).UTC()

	tr, err := NewTransaction(sym, source, dest, amount, memo, tm)
	assert.Nil(t, err)

	if assert.NotNil(t, tr) {
		assert.Equal(t, sym, tr.Symbol)
		assert.Equal(t, source, tr.Source)
		assert.Equal(t, dest, tr.Destination)
		assert.Equal(t, amount, tr.Amount)
		assert.Equal(t, memo, tr.Memo)
		assert.Equal(t, tm, tr.Time)
	}
}

// TestCalculateId verifies we're able to calculate a transaction hash.
func TestCalculateId(t *testing.T) {
	id := "ajRUhpCJAkyLKcsoXw8WTDR-VtEyftekhoFi4Air1LI+"
	sym := "TEST"
	source := "source_address"
	dest := "dest_address"
	amount := 1.0
	memo := "memo"
	tm := time.Unix(0, 0).UTC()

	tr, err := NewTransaction(sym, source, dest, amount, memo, tm)
	assert.Nil(t, err)

	err = tr.CalculateID()
	assert.Nil(t, err)
	assert.Equal(t, id, tr.ID)

	// ensure this returns the same value if called multiple times.
	err = tr.CalculateID()
	assert.Nil(t, err)
	assert.Equal(t, id, tr.ID)
}

//TestSignAndVerify tests whether we can sign transactions and verify their signature.
func TestSignAndVerify(t *testing.T) {
	k, err := keys.GenerateKeyPair()
	assert.Nil(t, err)

	source, err := keys.GetAddress(k)
	assert.Nil(t, err)

	sym := "TEST"
	dest := "dest_address"
	amount := 1.0
	memo := "memo"
	tm := time.Unix(0, 0).UTC()

	tr, err := NewTransaction(sym, source, dest, amount, memo, tm)
	assert.Nil(t, err)

	_, _, err = tr.SignTransaction(k)
	assert.Nil(t, err)

	ok, err := tr.VerifyTransaction()
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
}
