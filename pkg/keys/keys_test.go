package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateKeyPair checks if we're able to generate a keypair.
func TestGenerateKeyPair(t *testing.T) {
	k, err := GenerateKeyPair()

	assert.Nil(t, err)
	assert.NotNil(t, k)
}

// TestGetAddress confirms we're able to convert a public key into an address.
func TestGetAddress(t *testing.T) {
	k, err := GenerateKeyPair()
	assert.Nil(t, err)
	assert.NotNil(t, k)

	address, err := GetAddress(k)
	assert.Nil(t, err)
	assert.NotNil(t, address)
}

// TestDecodeAddress confirms we're able to convert an address key into a public key.
func TestDeodeAddress(t *testing.T) {
	k, err := GenerateKeyPair()
	assert.Nil(t, err)
	assert.NotNil(t, k)

	address, err := GetAddress(k)
	assert.Nil(t, err)
	assert.NotNil(t, address)

	pubK, err := DecodeAddress(address)
	assert.Nil(t, err)
	assert.NotNil(t, pubK)
	assert.Equal(t, k.PublicKey, *pubK)
}
