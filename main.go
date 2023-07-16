package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Shop struct {
	Name        string `json:"name"`
	WebhookURL  string `json:"webhookURL"`
	PublicKey   string `json:"publicKey"`
	Description string `json:"description"`
}

var db *sql.DB

func getShops(w http.ResponseWriter, r *http.Request) {
	log.Printf("Polling request received...Checking database....")
	rows, err := db.Query("SELECT name, webhookURL, publicKey, description FROM shops")
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		log.Printf("COMING HERE")
		var shop Shop
		var (
			shopnameVar string
			webhookURL  string
			PublicKey   string
			description string
		)

		if err := rows.Scan(&shopnameVar, &webhookURL, &PublicKey, &description); err != nil {
			fmt.Printf("Error! %s key is %s\n", shopnameVar, webhookURL, PublicKey, description)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		shop.Name = shopnameVar
		shop.WebhookURL = webhookURL
		shop.PublicKey = PublicKey
		shop.Description = description
		log.Printf("attempting append of %s to shops array", shop.Name)
		shops = append(shops, shop)
		log.Printf("appending successful")
	}

	log.Printf("appending loop finished")

	json.NewEncoder(w).Encode(shops)
}

func addShop(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	auth_server := os.Getenv("AUTH_SERVER")
	req, _ := http.NewRequest("GET", auth_server+"/validate", nil)
	req.Header.Add("Authorization", r.Header.Get("Authorization"))
	resp, err := client.Do(req)

	log.Printf("Checking authorization via Server")

	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var newShop Shop

	json.NewDecoder(r.Body).Decode(&newShop)

	log.Printf("Checking if shop already exists in database")

	err = db.QueryRow("SELECT name FROM shops WHERE name = $1", newShop.Name).Scan(&newShop.Name)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != sql.ErrNoRows {
		w.Write([]byte("You are already part of federation"))
		return
	}

	log.Printf("Inserting into database")

	_, err = db.Exec("INSERT INTO shops (name, webhookURL, publicKey, description) VALUES ($1, $2, $3, $4)", newShop.Name, newShop.WebhookURL, newShop.PublicKey, newShop.Description)
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
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	db, err = sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/shops", getShops).Methods("GET")
	router.HandleFunc("/shops", addShop).Methods("POST")

	port := os.Getenv("HUB_PORT")
	log.Printf("Federation hub is running on port %s", port)

	// Create a WaitGroup to wait for the server goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := http.ListenAndServe(":"+port, router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for an interrupt signal to cleanup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Signal the server goroutine to stop and wait for it to finish
	wg.Wait()

	db.Close()
	log.Println("Shutting down the server...")
}
