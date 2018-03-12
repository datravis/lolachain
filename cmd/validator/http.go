package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/datravis/lolachain/pkg/chain"
	"github.com/datravis/lolachain/pkg/tran"

	"github.com/gorilla/mux"
)

var lolachain *chain.Chain

// StartServer starts the validator HTTP server.
func StartServer(port string, c *chain.Chain) {
	lolachain = c
	r := mux.NewRouter()
	r.HandleFunc("/addresses/{address}", AddressHandler)
	r.HandleFunc("/transactions", TransactionHandler)
	r.HandleFunc("/chain", ChainHandler)
	r.HandleFunc("/pending", PendingHandler)
	http.Handle("/", r)

	fmt.Printf("%s", http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// ChainHandler returns the blockchain in JSON format.
func ChainHandler(w http.ResponseWriter, r *http.Request) {
	chainJSON, err := json.MarshalIndent(lolachain.Blocks, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(chainJSON))
}

func PendingHandler(w http.ResponseWriter, r *http.Request) {
	pendingJSON, err := json.MarshalIndent(lolachain.Pending, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(pendingJSON))
}

// AddressHandler returns balances for the supplied address.
func AddressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	balances := lolachain.GetBalanceForAddress(address)

	balancesJSON, err := json.Marshal(balances)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(balancesJSON))

}

// TransactionHandler handles posting new transactions to the blockchain.
func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	var t tran.Transaction
	err = json.Unmarshal(b, &t)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = lolachain.PostTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
