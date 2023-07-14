package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Shop struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhookURL"`
	PubKey_pem string `json:"pubKey_pem"`
}

type jwtToken struct {
	Token string `json:"token"`
}

var jwtKey = []byte("my_secret_key")
var refreshTokens = make(map[string]string)

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateJwtToken(name string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validateJwtToken(tokenString string) (bool, string) {
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, ""
	}

	if !token.Valid {
		return false, ""
	}

	return true, claims.Issuer
}

func login(w http.ResponseWriter, r *http.Request) {
	var shop Shop
	json.NewDecoder(r.Body).Decode(&shop)

	// Here you'd want to authenticate the shop, for example by checking the shop name against a database.
	// If the shop is authenticated, generate an access token for them.
	tokenString, err := generateJwtToken(shop.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshTokens[refreshToken] = shop.Name

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: tokenString, RefreshToken: refreshToken})
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	if originalName, ok := refreshTokens[body.RefreshToken]; ok {
		delete(refreshTokens, body.RefreshToken)
		accessToken, err := generateJwtToken(originalName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshToken, err := generateRefreshToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshTokens[refreshToken] = originalName

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{AccessToken: accessToken, RefreshToken: refreshToken})
	} else {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
	}
}

func validateToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	valid, _ := validateJwtToken(tokenString)

	if valid {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/token", refreshToken).Methods("POST")
	router.HandleFunc("/validate", validateToken).Methods("GET")

	log.Fatal(http.ListenAndServe(":8463", router))
}
