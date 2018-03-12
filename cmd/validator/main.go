package main

import (
	"flag"
	"fmt"

	"github.com/datravis/lolachain/pkg/chain"
	"github.com/datravis/lolachain/pkg/keys"
)

func main() {
	bind := flag.String("bind", "localhost:8081", "The address and port to bind the server to")
	seed := flag.String("seed", "", "A seed node to connect to")
	flag.Parse()

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

	c := &chain.Chain{
		Peers: []string{*seed},
	}

	go c.Validate(keyPair)

	StartServer(*bind, c)
}
