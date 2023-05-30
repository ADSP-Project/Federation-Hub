package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Shop struct {
	Name string `json:"name"`
}

var federationServer = "http://localhost:8000"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/webhook", handleWebhook).Methods("POST")

	go pollFederationServer()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	var newShop Shop
	json.NewDecoder(r.Body).Decode(&newShop)

	fmt.Printf("New shop joined the federation: %s\n", newShop.Name)
}

func pollFederationServer() {
	for {
		time.Sleep(10 * time.Second)

		resp, err := http.Get(federationServer + "/shops")
		if err != nil {
			log.Printf("Failed to poll federation server: %v\n", err)
			continue
		}

		var shops []Shop
		json.NewDecoder(resp.Body).Decode(&shops)

		fmt.Printf("Current shops in the federation: %v\n", shops)
	}
}
