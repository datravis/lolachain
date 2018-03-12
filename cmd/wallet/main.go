package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/datravis/lolachain/pkg/client"
	"github.com/datravis/lolachain/pkg/keys"
)

// TODO: This needs to be refactored. This is just a quickly thrown together
// implementation to get something going.
func main() {
	v := flag.String("validator", "http://localhost:8081", "the validator to connect to")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("At least one argument is required")
		return
	}

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

	address, err := keys.GetAddress(keyPair)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	switch command := args[0]; command {
	case "balance":
		if len(args) == 2 {
			address = args[1]
		}

		balances, err := client.GetBalances(*v, address)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		fmt.Printf("Address: %s\n", address)
		fmt.Println("Balances:")
		for key, val := range balances {
			fmt.Printf("%f %s\n", val, key)
		}
	case "send":
		if len(args) != 5 {
			fmt.Println("Requires arguments: dest amount symbol memo")
			return
		}
		dest := args[1]
		amount, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		symbol := args[3]
		memo := args[4]
		err = client.Send(*v, keyPair, dest, amount, symbol, memo)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	default:
		fmt.Println("Unknown command")
	}
}
