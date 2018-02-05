package main

import (
	"fmt"

	"github.com/datravis/lolachain/pkg/chain"
	"github.com/datravis/lolachain/pkg/keys"
)

func main() {
	c := &chain.Chain{}
	path, err := keys.GetDefaultKeyPath()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	keyPair, err := keys.LoadOrGenerateKeys(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	go c.Validate(keyPair)

	StartServer(c)
}
