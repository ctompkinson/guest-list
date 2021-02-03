package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"gorm.io/gorm"
	"net/http"
)

type getEmptySeatsResponse struct {
	SeatsEmpty int `json:"seats_empty"`
}

// HandleGetEmptySeats counts the amount of empty seats at the party right now
// it does not include guests that haven't checked in
func HandleGetEmptySeats(w http.ResponseWriter, r *http.Request) {
	db := database.Get()

	// Get all the tables and count their seats
	var tables []model.Table
	if err := db.Find(&tables).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to load tables: %v", err))
		return
	}
	totalSeats := 0
	for _, table := range tables {
		totalSeats = totalSeats + table.Seats
	}

	// Get all reservations and count used seats
	var reservations []model.Reservation
	if err := db.Where("arrival_time IS NOT NULL").Find(&reservations).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to load reservations: %v", err))
		return
	}
	usedSeats := 0
	for _, reservation := range reservations {
		usedSeats = usedSeats + reservation.AccompanyingGuests + 1 // Add on the primary guest
	}

	out, err := json.Marshal(getEmptySeatsResponse{SeatsEmpty: totalSeats - usedSeats})
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal response: %v", err))
		return
	}

	http.StatusText(http.StatusOK)
	_, _ = w.Write(out)
}
