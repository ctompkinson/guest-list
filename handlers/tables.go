package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type createTableRequest struct {
	Seats int `json:"seats"`
}

// HandleCreateTable creates a new table which can be used
// It must have a unique table number
func HandleCreateTable(w http.ResponseWriter, r *http.Request) {
	db := database.Get()
	tableNumber := mux.Vars(r)["tableNumber"]
	t, err := strconv.ParseInt(tableNumber, 10, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse table number: %v", err), http.StatusBadRequest)
		return
	}

	// Check if there is any tables with that number
	var table model.Table
	result := db.Where("number = ?", tableNumber).First(&table)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, fmt.Sprintf("failed to check if table exists: %v", result.Error), http.StatusInternalServerError)
		return
	}
	if result.RowsAffected != 0 {
		http.Error(w, "a table exists with that number already", http.StatusInternalServerError)
		return
	}

	// Get table seats
	var body createTableRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse body: %v", err), http.StatusBadRequest)
		return
	}

	// Create Table
	table = model.Table{
		Number: int(t),
		Seats:  body.Seats,
	}
	if err := db.Create(&table).Error; err != nil {
		http.Error(w, fmt.Sprintf("failed to create table: %v", err), http.StatusInternalServerError)
		return
	}

	http.StatusText(http.StatusOK)
	_, _ = w.Write([]byte(`{ "status": "created" }`))
}

// HandleDeleteTable deletes a table given its table number
func HandleDeleteTable(w http.ResponseWriter, r *http.Request) {
	db := database.Get()

	// Get the table number from the parameters and check it
	tableNumber := mux.Vars(r)["tableNumber"]
	if tableNumber == "" {
		http.Error(w, fmt.Sprintf("failed to give valid tableNumber: %v", tableNumber), http.StatusBadRequest)
		return
	}

	// Grab the table in question
	var table model.Table
	if err := db.Where("number = ?", tableNumber).First(&table).Error; err != nil {
		http.Error(w, fmt.Sprintf("unable to find table: %v", tableNumber), http.StatusInternalServerError)
		return
	}

	// Check if there is any reservations and stop table deletion
	var reservation model.Reservation
	result := db.Where("table_id = ?", table.ID).First(&reservation)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, fmt.Sprintf("failed to query reservations: %v", result.Error), http.StatusInternalServerError)
		return
	}
	if reservation.ID != 0 { // Same as a nil check
		http.Error(w, "cannot delete a table with a reservation", http.StatusInternalServerError)
		return
	}

	if err := db.Where("id = ?", table.ID).Unscoped().Delete(&model.Table{}).Error; err != nil {
		http.Error(w, fmt.Sprintf("failed to delete table: %v", err), http.StatusInternalServerError)
		return
	}

	http.StatusText(http.StatusOK)
	_, _ = w.Write([]byte(`{ "status": "deleted" }`))
}

// HandleGetTable gets the information about a table given its table number
func HandleGetTable(w http.ResponseWriter, r *http.Request) {
	db := database.Get()

	tableNumber := mux.Vars(r)["tableNumber"]
	if tableNumber == "" {
		http.Error(w, fmt.Sprintf("failed to give valid tableNumber: %v", tableNumber), http.StatusBadRequest)
		return
	}

	var t model.Table
	result := db.Where("number = ?", tableNumber).First(&t)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("failed to get table: %v", result.Error), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(t)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	http.StatusText(http.StatusOK)
	_, _ = w.Write(out)
}
