package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	router.HandleFunc("/orders", CreateOrder).Methods("POST")
	router.HandleFunc("/user/orders", GetUserOrders).Methods("GET")
	router.HandleFunc("/order/details", GetOrderDetails).Methods("GET")
	router.HandleFunc("/admin/orders", GetAllOrders).Methods("GET")
    router.HandleFunc("/admin/orders/update", UpdateOrderStatus).Methods("PUT")
    router.HandleFunc("/admin/orders/delete", DeleteOrder).Methods("DELETE")
	router.HandleFunc("/order/{order_id}/assign/{courier_id}", AssignOrderForCourier).Methods("POST")
	router.HandleFunc("/courier/{courier_id}/orders", GetCourierOrders).Methods("GET")
	router.HandleFunc("/order/{order_id}/decline/{courier_id}", DeclineOrder).Methods("DELETE")






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

	// Get the stored password, role, and user_id for the provided email
	var userID int
	var storedPassword, userRole string

	err = db.QueryRow("SELECT user_id, password, role FROM Users WHERE email = ?", credentials.Email).Scan(&userID, &storedPassword, &userRole)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found in the database
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"message": "User not found"}`, http.StatusUnauthorized)
		} else {
			// Some other error occurred
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	// Send JSON response with role and user_id
	response := map[string]interface{}{
		"message": "Login successful",
		"role":    userRole,
		"user_id":  userID,
	}
	json.NewEncoder(w).Encode(response)
}

type Order struct {
	UserID          int    `json:"user_id"`
	PickupLocation  string `json:"pickup_location"`
	DropoffLocation string `json:"dropoff_location"`
	PackageDetails  string `json:"package_details"`
	DeliveryTime    string `json:"delivery_time"`
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order

	// Parse JSON request body
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if order.UserID == 0 || order.PickupLocation == "" || order.DropoffLocation == "" || order.PackageDetails == "" || order.DeliveryTime == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Insert order into the database
	query := "INSERT INTO Orders (user_id, pickup_location, dropoff_location, package_details, delivery_time) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Exec(query, order.UserID, order.PickupLocation, order.DropoffLocation, order.PackageDetails, order.DeliveryTime)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Get the last inserted order ID
	orderID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve order ID", http.StatusInternalServerError)
		return
	}

	// Send success response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Order created successfully",
		"order_id":  orderID,
	})
}


func GetUserOrders(w http.ResponseWriter, r *http.Request) {
	// Parse user_id from the query parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Query to retrieve orders for the specific user, including status and courier_id
	rows, err := db.Query("SELECT order_id, pickup_location, dropoff_location, package_details, delivery_time, status, courier_id FROM Orders WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Define a struct to hold each order's details
	type Order struct {
		OrderID         int        `json:"order_id"`
		PickupLocation  string     `json:"pickup_location"`
		DropoffLocation string     `json:"dropoff_location"`
		PackageDetails  string     `json:"package_details"`
		DeliveryTime    string     `json:"delivery_time"` // DeliveryTime as string, no pointer needed
		Status          string     `json:"status"`
		CourierID       *int       `json:"courier_id"` // CourierID as pointer to handle NULL
	}

	// Slice to hold all orders
	var orders []Order

	// Iterate over query results and populate orders slice
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.OrderID,
			&order.PickupLocation,
			&order.DropoffLocation,
			&order.PackageDetails,
			&order.DeliveryTime, // Scan into string (no pointer)
			&order.Status,
			&order.CourierID, // Scan into pointer to handle NULL
		)
		if err != nil {
			http.Error(w, "Failed to parse order details: "+err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		return
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Respond with the list of orders in JSON format
	json.NewEncoder(w).Encode(orders)
}



func GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	var order struct {
		OrderID          int       `json:"order_id"`
		UserID           int       `json:"user_id"`
		PickupLocation   string    `json:"pickup_location"`
		DropoffLocation  string    `json:"dropoff_location"`
		PackageDetails   string    `json:"package_details"`
		DeliveryTime     *string   `json:"delivery_time,omitempty"`
		Status           string    `json:"status"`
		CreatedAt        *string   `json:"created_at,omitempty"`
	}

	var deliveryTimeRaw, createdAtRaw []byte

	query := `SELECT order_id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status, created_at 
              FROM Orders WHERE order_id = ?`

	err := db.QueryRow(query, orderID).Scan(&order.OrderID, &order.UserID, &order.PickupLocation, 
		&order.DropoffLocation, &order.PackageDetails, &deliveryTimeRaw, &order.Status, &createdAtRaw)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Convert the raw delivery time into a formatted string
	if len(deliveryTimeRaw) > 0 {
		// Attempt to parse delivery time
		deliveryTime, err := time.Parse("2006-01-02 15:04:05", string(deliveryTimeRaw))
		if err != nil {
			http.Error(w, "Error parsing delivery time: "+err.Error(), http.StatusInternalServerError)
			return
		}
		formattedTime := deliveryTime.Format(time.RFC3339)
		order.DeliveryTime = &formattedTime
	}

	// Convert the raw created_at into a formatted string
	if len(createdAtRaw) > 0 {
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtRaw))
		if err != nil {
			http.Error(w, "Error parsing created_at: "+err.Error(), http.StatusInternalServerError)
			return
		}
		formattedCreatedAt := createdAt.Format(time.RFC3339)
		order.CreatedAt = &formattedCreatedAt
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Prepare the response struct
	response := struct {
		OrderID          int       `json:"order_id"`
		UserID           int       `json:"user_id"`
		PickupLocation   string    `json:"pickup_location"`
		DropoffLocation  string    `json:"dropoff_location"`
		PackageDetails   string    `json:"package_details"`
		DeliveryTime     *string   `json:"delivery_time,omitempty"`
		Status           string    `json:"status"`
		CreatedAt        *string   `json:"created_at,omitempty"`
	}{
		OrderID:          order.OrderID,
		UserID:           order.UserID,
		PickupLocation:   order.PickupLocation,
		DropoffLocation:  order.DropoffLocation,
		PackageDetails:   order.PackageDetails,
		DeliveryTime:     order.DeliveryTime,
		Status:           order.Status,
		CreatedAt:        order.CreatedAt,
	}

	// Send JSON response
	json.NewEncoder(w).Encode(response)
}
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	// Prepare a slice to hold the orders
	var orders []struct {
		OrderID          int     `json:"order_id"`
		UserID           int     `json:"user_id"`
		PickupLocation   string  `json:"pickup_location"`
		DropoffLocation  string  `json:"dropoff_location"`
		PackageDetails   string  `json:"package_details"`
		DeliveryTime     string  `json:"delivery_time"`  // No *string, this field should always be filled
		Status           string  `json:"status"`
		CourierID        *int    `json:"courier_id,omitempty"` // Courier ID as a pointer (nullable)
	}

	// Query to retrieve all orders, excluding created_at
	query := `SELECT order_id, user_id, pickup_location, dropoff_location, package_details, 
                     delivery_time, status, courier_id FROM Orders`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to retrieve orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Loop through the rows and scan the data
	for rows.Next() {
		var order struct {
			OrderID          int     `json:"order_id"`
			UserID           int     `json:"user_id"`
			PickupLocation   string  `json:"pickup_location"`
			DropoffLocation  string  `json:"dropoff_location"`
			PackageDetails   string  `json:"package_details"`
			DeliveryTime     string  `json:"delivery_time"`  // No *string, this field should always be filled
			Status           string  `json:"status"`
			CourierID        *int    `json:"courier_id,omitempty"` // Courier ID as a pointer
		}

		// Scan the row into the order struct
		var deliveryTimeRaw []byte
		err := rows.Scan(&order.OrderID, &order.UserID, &order.PickupLocation, &order.DropoffLocation,
			&order.PackageDetails, &deliveryTimeRaw, &order.Status, &order.CourierID)
		if err != nil {
			http.Error(w, "Error scanning order data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert delivery_time if not NULL (must be non-NULL)
		if len(deliveryTimeRaw) == 0 {
			http.Error(w, "Order missing delivery time", http.StatusBadRequest)
			return
		}
		order.DeliveryTime = string(deliveryTimeRaw)

		// Append the order to the orders slice
		orders = append(orders, order)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send JSON response
	json.NewEncoder(w).Encode(orders)
}


func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the updated status
	var updatedStatus struct {
		Status string `json:"status"`
	}

	// Parse the request body to get the new status
	err := json.NewDecoder(r.Body).Decode(&updatedStatus)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Extract order ID from the URL (assuming it's passed as a URL parameter)
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	// Update the status in the Orders table
	updateOrderQuery := `UPDATE Orders SET status = ? WHERE order_id = ?`
	_, err = db.Exec(updateOrderQuery, updatedStatus.Status, orderID)
	if err != nil {
		http.Error(w, "Failed to update order status in Orders table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the status in the assignedorders table if there is an assigned record
	updateAssignedQuery := `UPDATE assignedorders SET status = ? WHERE order_id = ?`
	_, err = db.Exec(updateAssignedQuery, updatedStatus.Status, orderID)
	if err != nil {
		http.Error(w, "Failed to update order status in assignedorders table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order status updated successfully in both tables"})
}





func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order ID from the URL (assuming it's passed as a URL parameter)
	orderID := r.URL.Query().Get("order_id")

	// First, delete the record from the assignedorders table where the order_id matches
	deleteAssignedQuery := `DELETE FROM assignedorders WHERE order_id = ?`
	_, err := db.Exec(deleteAssignedQuery, orderID)
	if err != nil {
		http.Error(w, "Failed to delete from assignedorders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Then, delete the record from the orders table
	deleteOrderQuery := `DELETE FROM Orders WHERE order_id = ?`
	_, err = db.Exec(deleteOrderQuery, orderID)
	if err != nil {
		http.Error(w, "Failed to delete order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order and its assigned record deleted successfully"})
}



func AssignOrderForCourier(w http.ResponseWriter, r *http.Request) {
	// Parse the order_id and courier_id from the URL
	vars := mux.Vars(r)
	orderID := vars["order_id"]
	courierID := vars["courier_id"]

	// Get the status of the specific order from the orders table
	var orderStatus string
	err := db.QueryRow("SELECT status FROM Orders WHERE order_id = ?", orderID).Scan(&orderStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve order status: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if the order already exists in the assignedorders table
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM assignedorders WHERE order_id = ?)", orderID).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check assignedorders table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		// Update the courier_id in assignedorders if the order already exists
		_, err = db.Exec("UPDATE assignedorders SET courier_id = ?, status = ?, assigned_at = NOW() WHERE order_id = ?", courierID, orderStatus, orderID)
		if err != nil {
			http.Error(w, "Failed to update assignedorders: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Insert a new record in assignedorders if the order does not already exist
		_, err = db.Exec(`INSERT INTO assignedorders (order_id, courier_id, assigned_at, status) 
		                  VALUES (?, ?, NOW(), ?)`, orderID, courierID, orderStatus)
		if err != nil {
			http.Error(w, "Failed to insert into assignedorders: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Update the Orders table to assign the courier
	_, err = db.Exec("UPDATE Orders SET courier_id = ? WHERE order_id = ?", courierID, orderID)
	if err != nil {
		http.Error(w, "Failed to update order with courier ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order successfully assigned to courier"})
}



func GetCourierOrders(w http.ResponseWriter, r *http.Request) {
	// Parse the courier_id from the URL
	vars := mux.Vars(r)
	courierID := vars["courier_id"]

	// Prepare a slice to hold the orders assigned to the courier
	var orders []struct {
		OrderID          int    `json:"order_id"`
		UserID           int    `json:"user_id"`
		PickupLocation   string `json:"pickup_location"`
		DropoffLocation  string `json:"dropoff_location"`
		PackageDetails   string `json:"package_details"`
		DeliveryTime     string `json:"delivery_time"`
		Status           string `json:"status"`
		CourierID        int    `json:"courier_id"`
	}

	// Query to retrieve all orders assigned to the specific courier from the orders table
	query := `
		SELECT order_id, user_id, pickup_location, dropoff_location, 
		       package_details, delivery_time, status, courier_id
		FROM Orders
		WHERE courier_id = ?`

	// Execute the query
	rows, err := db.Query(query, courierID)
	if err != nil {
		http.Error(w, "Failed to retrieve orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Loop through the rows and scan the data into the orders slice
	for rows.Next() {
		var order struct {
			OrderID        int    `json:"order_id"`
			UserID         int    `json:"user_id"`
			PickupLocation string `json:"pickup_location"`
			DropoffLocation string `json:"dropoff_location"`
			PackageDetails string `json:"package_details"`
			DeliveryTime   string `json:"delivery_time"`
			Status         string `json:"status"`
			CourierID      int    `json:"courier_id"`
		}

		// Scan the row into the order struct
		err := rows.Scan(&order.OrderID, &order.UserID, &order.PickupLocation, &order.DropoffLocation,
			&order.PackageDetails, &order.DeliveryTime, &order.Status, &order.CourierID)
		if err != nil {
			http.Error(w, "Error scanning order data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Append the order to the orders slice
		orders = append(orders, order)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send JSON response
	json.NewEncoder(w).Encode(orders)
}


func DeclineOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order_id and courier_id from the URL path
	vars := mux.Vars(r)
	orderID := vars["order_id"]
	courierID := vars["courier_id"]

	// Check if both order_id and courier_id are provided
	if orderID == "" || courierID == "" {
		http.Error(w, "Both order_id and courier_id are required", http.StatusBadRequest)
		return
	}

	// Delete the record from the assignedorders table
	_, err := db.Exec("DELETE FROM assignedorders WHERE order_id = ? AND courier_id = ?", orderID, courierID)
	if err != nil {
		http.Error(w, "Failed to delete from assignedorders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the orders table to set courier_id to NULL for the given order_id
	_, err = db.Exec("UPDATE orders SET courier_id = NULL WHERE order_id = ?", orderID)
	if err != nil {
		http.Error(w, "Failed to update orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order declined successfully"})
}
