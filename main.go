package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Shop struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhookURL"`
	PubKey_pem string `json:"pubKey_pem"`
}

var db *sql.DB

func getShops(w http.ResponseWriter, r *http.Request) {
	log.Printf("Polling request received...Checking database....")
	rows, err := db.Query("SELECT name, pubkey_pem FROM shops")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		var shop Shop
		var (
			shopnameVar string
			pKeyVar     string
		)

		if err := rows.Scan(&shopnameVar, &pKeyVar); err != nil {
			fmt.Printf("Error! %s key is %s\n", shopnameVar, pKeyVar)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		shop.Name = shopnameVar
		shop.PubKey_pem = pKeyVar
		log.Printf("attempting append of %s with %s to shops array", shop.Name, shop.PubKey_pem)
		shops = append(shops, shop)
		log.Printf("appending successful")
	}

	log.Printf("appending loop finished")

	json.NewEncoder(w).Encode(shops)
}

func addShop(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/validate", nil)
	req.Header.Add("Authorization", r.Header.Get("Authorization"))
	resp, err := client.Do(req)

	log.Printf("Checking authorization via Server")

	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var newShop Shop

	json.NewDecoder(r.Body).Decode(&newShop)

	log.Printf("Inserting into database")

	_, err = db.Exec("INSERT INTO shops (name, webhookURL, pubKey_pem) VALUES ($1, $2, $3)", newShop.Name, newShop.WebhookURL, newShop.PubKey_pem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Select from Database")

	rows, err := db.Query("SELECT name, webhookURL FROM shops")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	log.Printf("Check database...")

	for rows.Next() {
		var shop Shop
		if err := rows.Scan(&shop.Name, &shop.WebhookURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go sendWebhook(shop.WebhookURL, newShop)

	}

	json.NewEncoder(w).Encode(newShop)

	log.Printf("Successfully added new shop")
}

func sendWebhook(webhookURL string, newShop Shop) {
	jsonData, _ := json.Marshal(newShop)

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send webhook to %s: %v\n", webhookURL, err)
		return
	}
	defer resp.Body.Close()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/shops", getShops).Methods("GET")
	router.HandleFunc("/shops", addShop).Methods("POST")

	port := ":8000"
	log.Printf("Federation hub is running on port%s", port)

	http.ListenAndServe(port, router)
}
