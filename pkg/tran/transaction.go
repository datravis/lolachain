package tran

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"time"

	"github.com/datravis/lolachain/pkg/keys"

	"github.com/lytics/base62"
)

// Transaction contains information about a transaction on the blockchain.
type Transaction struct {
	ID          string    `json:"id,omitempty"`
	Symbol      string    `json:"symbol"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Amount      float64   `json:"amount"`
	Memo        string    `json:"memo"`
	Time        time.Time `json:"time"`
	R           *big.Int  `json:"r,omitempty"`
	S           *big.Int  `json:"s,omitempty"`
}

// NewTransaction returns a new transaction.
func NewTransaction(symbol, source, dest string, amount float64, memo string, tm time.Time) (Transaction, error) {
	t := Transaction{
		Symbol:      symbol,
		Source:      source,
		Destination: dest,
		Amount:      amount,
		Memo:        memo,
		Time:        tm,
	}

	err := t.CalculateID()
	return t, err
}

// CalculateID calculates a transaction's ID.
func (t *Transaction) CalculateID() error {
	tmpTrans := Transaction{
		Symbol:      t.Symbol,
		Source:      t.Source,
		Destination: t.Destination,
		Amount:      t.Amount,
		Memo:        t.Memo,
		Time:        t.Time,
	}
	jsonEncoded, err := json.Marshal(tmpTrans)

	if err != nil {
		return err
	}

	hash := sha256.Sum256(jsonEncoded)
	t.ID = base62.StdEncoding.EncodeToString(hash[:])

	return nil

}

// SignTransaction sign's a transaction with a user's private key.
func (t *Transaction) SignTransaction(key *ecdsa.PrivateKey) (*big.Int, *big.Int, error) {
	tmpTrans := Transaction{
		ID:          t.ID,
		Symbol:      t.Symbol,
		Source:      t.Source,
		Destination: t.Destination,
		Amount:      t.Amount,
		Memo:        t.Memo,
		Time:        t.Time,
	}
	jsonEncoded, err := json.Marshal(tmpTrans)
	if err != nil {
		return nil, nil, err
	}

	sum := sha256.Sum256(jsonEncoded)
	r, s, err := ecdsa.Sign(rand.Reader, key, sum[:])
	t.R = r
	t.S = s

	return r, s, err
}

// VerifyTransaction verifies a transaction was signed by the proper private key.
func (t *Transaction) VerifyTransaction() (bool, error) {
	tmpTrans := Transaction{
		ID:          t.ID,
		Symbol:      t.Symbol,
		Source:      t.Source,
		Destination: t.Destination,
		Amount:      t.Amount,
		Memo:        t.Memo,
		Time:        t.Time,
	}
	jsonEncoded, err := json.Marshal(tmpTrans)
	if err != nil {
		return false, err
	}

	sum := sha256.Sum256(jsonEncoded)

	publicKey, err := keys.DecodeAddress(t.Source)
	if err != nil {
		return false, err
	}
	return ecdsa.Verify(publicKey, sum[:], t.R, t.S), nil
}
