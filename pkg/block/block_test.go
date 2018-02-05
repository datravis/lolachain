package block

import (
	"testing"
	"time"

	"github.com/datravis/lolachain/pkg/tran"

	"github.com/stretchr/testify/assert"
)

// TestNewBlock verifies we're able to create a new block using the function by the package.
func TestNewBlock(t *testing.T) {
	index := uint64(0)
	tm := time.Unix(0, 0).UTC()
	tr := []tran.Transaction{}
	address := "my_addr"
	previousHash := [32]byte{}

	b, err := NewBlock(index, tm, tr, address, previousHash)
	assert.Nil(t, err)

	if assert.NotNil(t, b) {
		assert.Equal(t, index, b.Index)
		assert.Equal(t, tm, b.Time)
		assert.Equal(t, tr, b.Transactions)
		assert.Equal(t, address, b.Validator)
		assert.Equal(t, previousHash, b.PreviousHash)
	}
}

// TestCalculateHash verifies the CalculateHash method returns the expected hash.
func TestCalculateHash(t *testing.T) {
	hash := [32]uint8{0x91, 0x46, 0xfb, 0xfb, 0x28, 0xaf, 0x8e, 0xbf, 0x8b, 0x7c, 0x71, 0x3e, 0x77, 0x6e, 0x22, 0x72, 0x13, 0x8c, 0x8a, 0xba, 0x92, 0x74, 0xb, 0x62, 0xee, 0xd1, 0x7c, 0x6b, 0x2f, 0x5d, 0x61, 0x64}
	index := uint64(0)
	tm := time.Unix(0, 0)
	tr := []tran.Transaction{}
	address := "my_addr"
	previousHash := [32]byte{}

	b, err := NewBlock(index, tm, tr, address, previousHash)
	assert.Nil(t, err)

	h, err := b.CalculateHash()
	assert.Nil(t, err)

	if assert.NotNil(t, h) {
		assert.Equal(t, hash, h)
		assert.Equal(t, hash, b.Hash)
	}
}
