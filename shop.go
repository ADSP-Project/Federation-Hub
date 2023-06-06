package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Shop struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhookURL"`
	PubKey_pem string `json:"pubKey_pem"`
}

var federationServer = "http://localhost:8000"

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run shop.go [port] [name]")
	}

	port := os.Args[1]
	shopName := os.Args[2]

	router := mux.NewRouter()
	router.HandleFunc("/webhook", handleWebhook).Methods("POST")

	go joinFederation(shopName)
	go pollFederationServer()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	var newShop Shop
	json.NewDecoder(r.Body).Decode(&newShop)

	fmt.Printf("New shop joined the federation: %s\n", newShop.Name)

	fmt.Printf("Public Key: %s", newShop.PubKey_pem)
}

func ExportPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {
	pubkey_pem := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
	return pubkey_pem
}
func ExportPrivateKeyAsPemStr(privatekey *rsa.PrivateKey) string {
	privatekey_pem := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}))
	return privatekey_pem
}

func joinFederation(shopName string) {

	PrivKey, err := rsa.GenerateKey(rand.Reader, 128)
	privatekey_pem := ExportPrivateKeyAsPemStr(PrivKey)
	pubKey_pem := ExportPublicKeyAsPemStr(&PrivKey.PublicKey)

	newShop := Shop{Name: shopName, WebhookURL: fmt.Sprintf("http://localhost:%s/webhook", os.Args[1]), PubKey_pem: pubKey_pem}

	log.Printf("New Shop Private Key is %s", privatekey_pem)
	log.Printf("New Shop Public key is %s", newShop.PubKey_pem)

	resp, err := http.PostForm("http://localhost:8081/login", url.Values{"name": {shopName}, "webhookURL": {newShop.WebhookURL}, "pubKey_pem": {newShop.PubKey_pem}})
	if err != nil {
		log.Fatal("Failed to authenticate with auth server")
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	accessToken := result["access_token"]

	jsonData, _ := json.Marshal(newShop)
	req, err := http.NewRequest("POST", federationServer+"/shops", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", accessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Failed to join federation: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Shop joined the federation")
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
