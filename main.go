package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Shop struct {
	Name string `json:"name"`
}

var shopList = struct{
    sync.RWMutex
    shops []Shop
}{}

func getShops(w http.ResponseWriter, r *http.Request) {
    shopList.RLock()
    defer shopList.RUnlock()
	json.NewEncoder(w).Encode(shopList.shops)
}

func addShop(w http.ResponseWriter, r *http.Request) {
	var newShop Shop
	json.NewDecoder(r.Body).Decode(&newShop)
    shopList.Lock()
    shopList.shops = append(shopList.shops, newShop)
    shopList.Unlock()
	json.NewEncoder(w).Encode(newShop)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/shops", getShops).Methods("GET")
	router.HandleFunc("/shops", addShop).Methods("POST")
	http.ListenAndServe(":8000", router)
}
