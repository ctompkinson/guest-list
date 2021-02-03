package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"net/http"
	"time"
)

type guestArrivalRequest struct {
	AccompanyingGuests int `json:"accompanying_guests"`
}
type guestArrivalResponse struct {
	Name string `json:"name"`
}

// HandleGuestArrival lets you signal that a guest has arrived at the party given a guests name and
// the amount of guests they have shown up with
func HandleGuestArrival(w http.ResponseWriter, r *http.Request) {
	db := database.Get()
	guestName := mux.Vars(r)["name"]
	if guestName == "" {
		ErrorResponse(w, http.StatusBadRequest, "unable to retrieve guest name from URL")
		return
	}

	// Decode body into request struct
	var reqBody guestArrivalRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("unable to parse body: %v", err))
		return
	}

	// Find the reservation
	var reservation model.Reservation
	if err := db.Preload("Table").Where("guest = ?", guestName).First(&reservation).Error; err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("unable to find guest: %v", err))
		return
	}

	// Does our reservation have a different amount of seats?
	if reservation.AccompanyingGuests != reqBody.AccompanyingGuests {
		// Were going to check our own table by checking reservations, so we only want to check for new guests (minus overselves)
		newGuests := reqBody.AccompanyingGuests - reservation.AccompanyingGuests

		enoughSeats, err := areEnoughSeatsAvailable(reservation.Table, newGuests)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to check available seats: %v", err))
			return
		}
		if !enoughSeats {
			ErrorResponse(w, http.StatusInternalServerError, "not enough seats available on selected table")
			return
		}
	}

	// We can now update our reservation and add an arrival time
	reservation.AccompanyingGuests = reqBody.AccompanyingGuests
	now := time.Now()
	reservation.ArrivalTime = &now

	if err := db.Save(&reservation).Error; err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to save update to reservation: %v", err))
		return
	}

	out, err := json.Marshal(guestArrivalResponse{Name: guestName})
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal response: %v", err))
		return
	}
	http.StatusText(http.StatusOK)
	_, _ = w.Write(out)
}

// HandleListGuests lists guests that have arrived at the party and their arrival time
func HandleListGuests(w http.ResponseWriter, r *http.Request) {
	db := database.Get()

	var reservations []model.Reservation
	if err := db.Preload("Table").Where("arrival_time IS NOT NULL").Find(&reservations).Error; err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to load reservations: %v", err))
		return
	}

	formattedReservations := []model.FormattedGuestArrival{}
	for _, res := range reservations {
		formattedReservations = append(formattedReservations, res.FormatAsGuestArrival())
	}

	out, err := json.Marshal(map[string][]model.FormattedGuestArrival{
		"guests": formattedReservations,
	})
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal reservations: %v", err))
		return
	}

	http.StatusText(http.StatusOK)
	_, _ = w.Write(out)
}
