package chain

import (
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"time"

	"github.com/datravis/lolachain/pkg/block"
	"github.com/datravis/lolachain/pkg/client"
	"github.com/datravis/lolachain/pkg/keys"
	"github.com/datravis/lolachain/pkg/tran"
)

const INCREMENTOR_DIVISOR = 128457181

// Chain contains a chain of blocks a long with pending transactions.
type Chain struct {
	Blocks    []*block.Block
	Pending   []tran.Transaction
	Peers     map[string]bool
	MyAddress string
}

// NewChain returns an instance of a chain.
func NewChain(address string, peers map[string]bool) *Chain {
	return &Chain{
		Blocks:    make([]*block.Block, 0, 0),
		Pending:   []tran.Transaction{},
		Peers:     peers,
		MyAddress: address,
	}
}

// Validate runs the main validator processing loop.
func (c *Chain) Validate(keyPair *ecdsa.PrivateKey) {
	c.Blocks = c.FetchBlocks()
	fmt.Printf("Fetched %d blocks from peer\n", len(c.Blocks))

	if len(c.Blocks) == 0 {
		_, err := c.GenesisBlock(keyPair)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			c.NotifyPeers()
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			c.FindPeers()
		}
	}()

	for {
		done := make(chan interface{})
		incrementorStream := c.FindIncrementor(done)
		blockUpdateStream := c.FindBlockUpdates(done)
		select {
		case count := <-blockUpdateStream:
			fmt.Printf("Fetched %d blocks from peer\n", count)
			fmt.Printf("Chain length: %d\n", len(c.Blocks))
		case inc := <-incrementorStream:
			blk, err := c.NextBlock(c.Pending, inc, keyPair)

			if err != nil {
				fmt.Printf("Error: %s\n", err)
			} else {
				fmt.Printf("New Block Generated: %d\n", blk.Index)
				fmt.Printf("Chain length: %d\n", len(c.Blocks))
				c.Pending = make([]tran.Transaction, 0)
			}

		}
		close(done)
	}
}

func (c *Chain) FindPeers() {
	for peer, _ := range c.Peers {
		peers, err := client.GetPeers(peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		for p, _ := range peers {
			c.Peers[p] = true
		}
	}
}

// PostTransaction posts a transaction the transaction mempool.
func (c *Chain) PostTransaction(t tran.Transaction) error {
	ok, err := t.VerifyTransaction()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Transaction invalid: %s", t.ID)
	}

	ok, err = c.VerifyBalance(t)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Transaction invalid: %s", t.ID)
	}

	c.Pending = append(c.Pending, t)
	return nil
}

// NextBlock Adds the next block to the blockchain.
func (c *Chain) NextBlock(transactions []tran.Transaction, incrementor uint64, keyPair *ecdsa.PrivateKey) (*block.Block, error) {
	lastBlock := c.Blocks[len(c.Blocks)-1]
	ts := time.Now().UTC()

	reward, err := c.CreateRewardTransaction(ts, "RKY", keyPair)
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, reward)

	reward, err = c.CreateRewardTransaction(ts, "LOLA", keyPair)
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, reward)

	validTransactions := c.ValidateTransactions(transactions)

	validatorAddress, err := keys.GetAddress(keyPair)
	if err != nil {
		return nil, err
	}

	nextBlock, err := block.NewBlock(lastBlock.Index+1, ts, validTransactions, validatorAddress, lastBlock.Hash, incrementor)
	if err != nil {
		return nil, err
	}

	c.Blocks = append(c.Blocks, nextBlock)
	return nextBlock, nil
}

// GenesisBlock produces the initial block of the blockchain.
func (c *Chain) GenesisBlock(keyPair *ecdsa.PrivateKey) (*block.Block, error) {
	validatorAddress, err := keys.GetAddress(keyPair)
	if err != nil {
		return nil, err
	}

	ts := time.Now().UTC()
	nextBlock, err := block.NewBlock(0, ts, []tran.Transaction{}, validatorAddress, [32]byte{}, INCREMENTOR_DIVISOR)
	if err != nil {
		return nil, err
	}

	c.Blocks = append(c.Blocks, nextBlock)
	return nextBlock, nil
}

// ValidateTransactions validates the list of provided transactions.
func (c *Chain) ValidateTransactions(trans []tran.Transaction) []tran.Transaction {
	batch := []tran.Transaction{}
	for _, t := range trans {
		if t.Memo != "block reward" {
			ok, err := c.VerifyBalance(t)
			if err != nil || !ok {
				fmt.Printf("transaction invalid: %s\n", err)
				continue
			}
		}
		ok, err := t.VerifyTransaction()
		if err != nil || !ok {
			fmt.Printf("transaction invalid: %s\n", err)
			continue
		}

		batch = append(batch, t)
	}

	return batch
}

// CreateRewardTransaction returns a block reward transaction.
func (c *Chain) CreateRewardTransaction(ts time.Time, symbol string, keyPair *ecdsa.PrivateKey) (tran.Transaction, error) {
	validatorAddress, err := keys.GetAddress(keyPair)
	if err != nil {
		return tran.Transaction{}, err
	}
	t, err := tran.NewTransaction(symbol, validatorAddress, validatorAddress, 1, "block reward", ts)
	if err != nil {
		return tran.Transaction{}, err
	}

	_, _, err = t.SignTransaction(keyPair)

	return t, err
}

// GetBalanceForAddress computes the balance of an address.
func (c *Chain) GetBalanceForAddress(a string) map[string]float64 {
	balances := make(map[string]float64)

	for _, b := range c.Blocks {
		for _, t := range b.Transactions {
			if t.Source == a || t.Destination == a {
				if _, ok := balances[t.Symbol]; !ok {
					balances[t.Symbol] = 0.0
				}
			}
			if t.Source == a && t.Destination == a && t.Memo == "block reward" {
				balances[t.Symbol] += t.Amount
			} else if t.Source == a {
				balances[t.Symbol] -= t.Amount
			} else if t.Destination == a {
				balances[t.Symbol] += t.Amount
			}
		}
	}

	return balances
}

// VerifyBalance confirms that a wallet contains a sufficient balance.
func (c *Chain) VerifyBalance(t tran.Transaction) (bool, error) {
	cBalance := c.GetBalanceForAddress(t.Source)
	if _, ok := cBalance[t.Symbol]; !ok {
		return false, fmt.Errorf("Insufficient funds to perform transaction")
	}
	if cBalance[t.Symbol] < t.Amount {
		return false, fmt.Errorf("Insufficient funds to perform transaction")
	}

	return true, nil
}

// FindIncrementor implements a simple proof of work algorithm.
func (c *Chain) FindIncrementor(done chan interface{}) <-chan uint64 {
	fmt.Println("Finding next incrementor")
	incrementorStream := make(chan uint64)
	go func() {
		defer close(incrementorStream)
		incrementor := rand.Uint64()
		for {

			select {
			case <-done:
				return
			default:
				if incrementor%INCREMENTOR_DIVISOR == 0 && incrementor != 0 {
					incrementorStream <- incrementor
					return
				}
				incrementor = rand.Uint64()
			}
		}
	}()

	return incrementorStream
}

// FindBlockUpdates retrieves the longest blockchain from our peers.
func (c *Chain) FindBlockUpdates(done chan interface{}) <-chan int {
	blockUpdateStream := make(chan int)
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		defer close(blockUpdateStream)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				longestChain := c.FetchBlocks()
				diff := len(longestChain) - len(c.Blocks)
				if diff > 0 {
					c.Blocks = longestChain
					blockUpdateStream <- diff
					return
				}
			}
		}
	}()

	return blockUpdateStream
}

// FetchBlocks fetches the longest blockchain from our peers.
func (c *Chain) FetchBlocks() []*block.Block {
	blocks := make([]*block.Block, 0, 0)
	for peer, _ := range c.Peers {
		tmpBlocks, err := client.GetBlocks(peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		if len(tmpBlocks) > len(blocks) {
			blocks = tmpBlocks
		}
	}

	return blocks
}

// NotifyPeers notifies our peers of our existance.
func (c *Chain) NotifyPeers() {
	for peer, _ := range c.Peers {
		client.PostPeer(peer, c.MyAddress)
	}
}

// AddPeer adds a new peer to our collection of peers.
func (c *Chain) AddPeer(peer string) {
	c.Peers[peer] = true
}
