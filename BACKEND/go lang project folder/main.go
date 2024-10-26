package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors" // Import the CORS package
	//"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

func main() {
	// Database connection
	db, err = sql.Open("mysql", "root:omar1234@tcp(127.0.0.1:3306)/package_tracking_system")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize router
	router := mux.NewRouter()

	// Route Handlers
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Allow Angular app only
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start the server with CORS
	fmt.Println("Server starting on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", c.Handler(router))) // Wrap the router with the CORS handler
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	// Parse JSON request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert user into the database without password hashing
	query := "INSERT INTO Users (name, email, phone, password, role) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(query, user.Name, user.Email, user.Phone, user.Password, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Send JSON response
	response := map[string]string{"message": "User registered successfully"}
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON request body
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the stored password for the provided email
	var storedPassword string
	err = db.QueryRow("SELECT password FROM Users WHERE email = ?", credentials.Email).Scan(&storedPassword)
	if err != nil {
		// User not found
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"message": "User not found"}`, http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored password
	if credentials.Password != storedPassword {
		// Invalid password
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"message": "Invalid password"}`, http.StatusUnauthorized)
		return
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send JSON response
	response := map[string]string{"message": "Login successful"}
	json.NewEncoder(w).Encode(response)
}

