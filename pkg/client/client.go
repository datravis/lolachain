package client

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/datravis/lolachain/pkg/block"
	"github.com/datravis/lolachain/pkg/keys"
	"github.com/datravis/lolachain/pkg/tran"
)

// GetBalances return's a wallet's balances.
func GetBalances(host string, address string) (map[string]float64, error) {
	balances := make(map[string]float64)

	url := fmt.Sprintf("%s/addresses/%s", host, address)
	resp, err := http.Get(url)
	if err != nil {
		return balances, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return balances, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &balances)
	return balances, err
}

// Send submits a new transaction to the lolachain API.
func Send(host string, keyPair *ecdsa.PrivateKey, dest string, amount float64, symbol string, memo string) error {
	address, err := keys.GetAddress(keyPair)
	if err != nil {
		return err
	}

	ts := time.Now().UTC()
	t, err := tran.NewTransaction(symbol, address, dest, amount, memo, ts)
	if err != nil {
		return err
	}
	_, _, err = t.SignTransaction(keyPair)
	if err != nil {
		return err
	}

	tJSON, err := json.Marshal(t)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/transactions", host), "application/json", bytes.NewBuffer(tJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return errors.New(string(body))
}

// GetBalances return's validators block chain.
func GetBlocks(host string) ([]*block.Block, error) {
	blocks := make([]*block.Block, 0)

	url := fmt.Sprintf("%s/chain", host)
	resp, err := http.Get(url)
	if err != nil {
		return blocks, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return blocks, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &blocks)
	return blocks, err
}
