package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/datravis/lolachain/pkg/client"
	"github.com/datravis/lolachain/pkg/keys"
)

// PageVariables contains variables returned to the screen.
type PageVariables struct {
	RKYBalance  float64
	LOLABalance float64
	Address     string
}

//RKYWallet handles GET and POST requests.
func RKYWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ShowWallet(w, r)
	}

	if r.Method == "POST" {
		Send(w, r)
		ShowWallet(w, r)
	}
}

// ShowWallet shows a wallet's information.
func ShowWallet(w http.ResponseWriter, r *http.Request) {
	path, err := keys.GetDefaultKeyPath()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	keyPair, err := keys.LoadOrGenerateKeys(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	address, err := keys.GetAddress(keyPair)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	balances, err := client.GetBalances(validator, address)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rkyBalance := 0.0
	if _, ok := balances["RKY"]; ok {
		rkyBalance = balances["RKY"]
	}

	llaBalance := 0.0
	if _, ok := balances["LOLA"]; ok {
		llaBalance = balances["LOLA"]
	}

	WalletVars := PageVariables{
		RKYBalance:  rkyBalance,
		LOLABalance: llaBalance,
		Address:     address,
	}

	t, err := template.ParseFiles("templates/wallet.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = t.Execute(w, WalletVars)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Send submits a new transaction to the lolachain API.
func Send(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	form := r.Form
	dest := form.Get("destination")
	amountStr := form.Get("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	symbol := form.Get("symbol")
	memo := form.Get("memo")

	path, err := keys.GetDefaultKeyPath()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	keyPair, err := keys.LoadOrGenerateKeys(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = client.Send(validator, keyPair, dest, amount, symbol, memo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
