package block

import (
	"crypto/sha256"
	"encoding/json"
	"strconv"
	"time"

	"github.com/datravis/lolachain/pkg/tran"
)

// Block is an entry on the blockchain.
type Block struct {
	Index        uint64             `json:"index"`
	Time         time.Time          `json:"time"`
	Transactions []tran.Transaction `json:"transactions"`
	Validator    string             `json:"validator"`
	PreviousHash [32]byte           `json:"previous_hash"`
	Hash         [32]byte           `json:"hash"`
	Incrementor  uint64             `json:"incrementor"`
}

// NewBlock returns an instance of a Block based on the supplied parameters.
func NewBlock(index uint64, t time.Time, transactions []tran.Transaction, validator string, previousHash [32]byte, incrementor uint64) (*Block, error) {
	block := &Block{
		Index:        index,
		Time:         t,
		Transactions: transactions,
		Validator:    validator,
		PreviousHash: previousHash,
		Incrementor:  incrementor,
	}

	var err error
	block.Hash, err = block.CalculateHash()
	return block, err
}

// CalculateHash computes a blocks hash.
func (b *Block) CalculateHash() ([32]byte, error) {
	indexBytes := []byte(strconv.FormatUint(b.Index, 10))
	timeBytes := []byte(b.Time.UTC().Format(time.RFC3339))
	incrementorBytes := []byte(strconv.FormatUint(b.Incrementor, 10))

	transjson, err := json.Marshal(b.Transactions)
	if err != nil {
		return [32]byte{}, err
	}

	toHash := []byte{}
	toHash = append(toHash, indexBytes...)
	toHash = append(toHash, timeBytes...)
	toHash = append(toHash, transjson...)
	toHash = append(toHash, b.PreviousHash[:]...)
	toHash = append(toHash, incrementorBytes...)

	return sha256.Sum256(toHash), nil
}
