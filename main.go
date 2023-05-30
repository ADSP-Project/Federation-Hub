package main

import (
	"fmt"
	"log"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Shop struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhookURL"`
}

var db *sql.DB

func getShops(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT name FROM shops")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		var shop Shop
		if err := rows.Scan(&shop.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		shops = append(shops, shop)
	}

	json.NewEncoder(w).Encode(shops)
}

func addShop(w http.ResponseWriter, r *http.Request) {
	var newShop Shop
	json.NewDecoder(r.Body).Decode(&newShop)

	_, err := db.Exec("INSERT INTO shops (name, webhookURL) VALUES ($1, $2)", newShop.Name, newShop.WebhookURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT name, webhookURL FROM shops")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var shop Shop
		if err := rows.Scan(&shop.Name, &shop.WebhookURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go sendWebhook(shop.WebhookURL, newShop)
	}

	json.NewEncoder(w).Encode(newShop)
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
