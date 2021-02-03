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
)

type createGuestListRequest struct {
	TableNumber        int `json:"table"`
	AccompanyingGuests int `json:"accompanying_guests"`
}
type createGuestListResponse struct {
	Name string `json:"name"`
}

// HandleCreateReservation creates a new reservation given a primary guest,
// the amount of guests and a valid table number
func HandleCreateReservation(w http.ResponseWriter, r *http.Request) {
	db := database.Get()
	// POST /guest_list/{name}
	// { "table": int, "accompanying_guests": int }

	// Get guest name from URL params
	guestName := mux.Vars(r)["name"]
	if guestName == "" {
		ErrorResponse(w, http.StatusBadRequest, "unable to retrieve guest name from URL")
		return
	}

	// Decode body into request struct
	var reqBody createGuestListRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("unable to parse body: %v", err))
		return
	}

	// Find the table
	var table model.Table
	if err := db.Where("number = ?", reqBody.TableNumber).First(&table).Error; err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to find table: %v", err))
		return
	}

	// Check if that person has any reservations already under their name
	var existingGuestReservations []model.Reservation
	result := db.Where("guest = ?", guestName).First(&existingGuestReservations)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to query for existing guest reservations: %v", result.Error))
		return
	}
	if result.Error == nil {
		ErrorResponse(w, http.StatusInternalServerError, "the guest already has a reservation")
		return
	}

	// Check for all our guests, plus the main guest
	enoughSeats, err := areEnoughSeatsAvailable(table, reqBody.AccompanyingGuests+1)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to check available seats: %v", err))
		return
	}
	if !enoughSeats {
		ErrorResponse(w, http.StatusInternalServerError, "not enough seats available on selected table")
		return
	}

	// Now we can finally create the reservation
	reservation := model.Reservation{
		Guest:              guestName,
		AccompanyingGuests: reqBody.AccompanyingGuests,
		Table:              table,
	}
	if err := db.Create(&reservation).Error; err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create reservations: %v", err))
		return
	}

	out, err := json.Marshal(createGuestListResponse{Name: guestName})
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal response: %v", err))
		return
	}
	http.StatusText(http.StatusOK)
	_, _ = w.Write(out)
}

// HandleDeleteReservation deletes a reservation given the primary guests name
func HandleDeleteReservation(w http.ResponseWriter, r *http.Request) {
	db := database.Get()

	guestName := mux.Vars(r)["name"]
	if guestName == "" {
		http.Error(w, "unable to retrieve guest name from URL", http.StatusBadRequest)
		return
	}

	// Check if reservation exists
	var reservation model.Reservation
	if err := db.Where("guest = ?", guestName).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "guest does not have a reservation", http.StatusBadRequest)
			return
		}

		http.Error(w, fmt.Sprintf("failed to lookup reservation: %v", err), http.StatusInternalServerError)
		return
	}

	// Delete table, use unscoped to ensure a hard delete and not soft
	if err := db.Unscoped().Delete(&reservation).Error; err != nil {
		http.Error(w, fmt.Sprintf("failed to delete reservation: %v", err), http.StatusInternalServerError)
		return
	}

	http.StatusText(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"deleted"}`))
}

// HandleGetReservations gets all the existing reservations
func HandleGetReservations(w http.ResponseWriter, r *http.Request) {
	db := database.Get()

	var reservations []model.Reservation
	if err := db.Preload("Table").Find(&reservations).Error; err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to load reservations: %v", err))
		return
	}

	formattedReservations := []model.FormattedReservation{}
	for _, res := range reservations {
		formattedReservations = append(formattedReservations, res.FormatAsReservation())
	}

	out, err := json.Marshal(map[string][]model.FormattedReservation{
		"guests": formattedReservations,
	})
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal reservations: %v", err))
		return
	}

	http.StatusText(http.StatusOK)
	_, _ = w.Write(out)
}
