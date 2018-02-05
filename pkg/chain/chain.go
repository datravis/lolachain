package chain

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/datravis/lolachain/pkg/block"
	"github.com/datravis/lolachain/pkg/keys"
	"github.com/datravis/lolachain/pkg/tran"
)

// Chain contains a chain of blocks a long with pending transactions.
type Chain struct {
	Blocks  []*block.Block
	Pending []tran.Transaction
}

// Validate runs the main validator processing loop.
func (c *Chain) Validate(keyPair *ecdsa.PrivateKey) {
	_, err := c.GenesisBlock(keyPair)
	if err != nil {
		fmt.Println(err.Error())
	}

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			blk, err := c.NextBlock(c.Pending, keyPair)

			if err != nil {
				fmt.Printf("Error: %s\n", err)
			} else {
				fmt.Printf("New Block Generated: %d\n", blk.Index)
				for _, t := range blk.Transactions {
					fmt.Printf("tx:\n src: %s\n dest: %s\n amt: %f %s\n memo: %s\n\n", t.Source, t.Destination, t.Amount, t.Symbol, t.Memo)
				}
				c.Pending = make([]tran.Transaction, 0)
			}

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
func (c *Chain) NextBlock(transactions []tran.Transaction, keyPair *ecdsa.PrivateKey) (*block.Block, error) {
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

	nextBlock, err := block.NewBlock(lastBlock.Index+1, ts, validTransactions, validatorAddress, lastBlock.Hash)
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
	nextBlock, err := block.NewBlock(0, ts, []tran.Transaction{}, validatorAddress, [32]byte{})
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
