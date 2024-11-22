package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/google/uuid"
)

// MODEL part
type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        string  `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type IDResponse struct {
	ID string `json:"id"`
}


// CONTROLLER part
var receipts = make(map[string]Receipt) // init an empty list of receipts

// Calculate the points by the rules
func calculatePoints(receipt Receipt) int {
	points := 0

	// 1. One point for every alphanumeric character in the retailer name
	points += len(regexp.MustCompile(`\w`).FindAllString(receipt.Retailer, -1))

	// 2. 50 points if the total is a round dollar amount with no cents
	if total, err := parseFloat(receipt.Total); err == nil && total == float64(int(total)) {
		points += 50
	}

	// 3. 25 points if the total is a multiple of 0.25
	if total, err := parseFloat(receipt.Total); err == nil && math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// 4. 5 points for every two items on the receipt
	points += 5 * (len(receipt.Items) / 2)

	// 5. If the trimmed length of the item description is a multiple of 3, apply the rule
	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			if price, err := parseFloat(item.Price); err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// 6. 6 points if the day in the purchase date is odd
	if date, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil && date.Day()%2 != 0 {
		points += 6
	}

	// 7. 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
			points += 10
		}
	}

	return points
}

// Check the total is float or not
func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// Add new receipt in the server
func addReceipt(w http.ResponseWriter, r *http.Request) {
	// Parse the receipt into JSON style
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil || !isValidReceipt(receipt) {
		http.Error(w, "Invalid receipt", http.StatusBadRequest)
		return
	}

	// Generate ID for the receipt and add it to the list
	receiptID := uuid.New().String()
	receipts[receiptID] = receipt

	// Return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := IDResponse{
		ID: receiptID,
	}
	json.NewEncoder(w).Encode(response)
}

// Get receipt points
func getReceiptPoints(w http.ResponseWriter, r *http.Request) {
	// Extract the receipt ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL structure", http.StatusBadRequest)
		return
	}

	id := pathParts[2]

	// Check if the receipt exists
	receipt, found := receipts[id]
	if !found {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Calculate points based on the provided rules
	points := calculatePoints(receipt)

	// Return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"points": int64(points)})
}

// Self-add fucntion to see if the receipts list is correct
func getAllReceipts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipts)
}

// Check the receipt is valid or not by retailer and total rule
func isValidReceipt(receipt Receipt) bool {
	// Checking retailer
	if !regexp.MustCompile(`^[\w\s\-\&]+$`).MatchString(receipt.Retailer) {
		return false
	}
	// Checking total
	if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(receipt.Total) {
		return false
	}
	
	return true
}

func main() {
	// Define the api routes
	http.HandleFunc("/receipts/process", addReceipt)     // Add new receipt, will return the receipt ID
	http.HandleFunc("/receipts/", getReceiptPoints)      // Get receipt points by receipt ID
	http.HandleFunc("/receipts/all", getAllReceipts)     // Get all receipts

	// Start the server using port 8080 and listen to event
	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(port, nil)
}
