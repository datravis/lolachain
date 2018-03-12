package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	validator string
)

// TODO: This should be rewritten to not be a client/server model.
func main() {
	port := flag.String("port", "8080", "The port to bind the server to")
	v := flag.String("validator", "http://localhost:8081", "the validator to connect to")
	flag.Parse()

	validator = *v

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.HandleFunc("/", RKYWallet)

	err := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
