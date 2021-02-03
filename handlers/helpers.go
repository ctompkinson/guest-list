package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"gorm.io/gorm"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

// ErrorResponse generates a json response given a http status code and a message
func ErrorResponse(w http.ResponseWriter, code int, message string) {
	res, _ := json.Marshal(errorResponse{Message: message})
	http.Error(w, string(res), code)
}

// areEnoughSeatsAvailable checks if enough seats are available on a table
// ensure that newGuests is the exact amount of people who you want to put on the table
func areEnoughSeatsAvailable(table model.Table, newGuests int) (bool, error) {
	db := database.Get()

	// Lets all reservations for our table
	var existingTableReservations []model.Reservation
	result := db.Where("table_id = ?", table.ID).Find(&existingTableReservations)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, result.Error
	}

	seatsUsed := 0
	if len(existingTableReservations) > 0 {
		for _, r := range existingTableReservations {
			seatsUsed = seatsUsed + r.AccompanyingGuests + 1 // Plus one to account the primary guest
		}
	}
	if (table.Seats - seatsUsed) >= newGuests {
		return true, nil
	}
	return false, nil
}
