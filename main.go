package main

import (
	"fmt"
	"log"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Shop struct {
	Name string `json:"name"`
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

	_, err := db.Exec("INSERT INTO shops (name) VALUES ($1)", newShop.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newShop)
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
	http.ListenAndServe(":8000", router)
}
