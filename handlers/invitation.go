package handlers

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"github.com/ctompkinson/guest-list/templates"
	"html/template"
	"net/http"
)

// HandleCreateInvitation creates a HTML invitation for a given guest
func HandleCreateInvitation(w http.ResponseWriter, r *http.Request) {
	db := database.Get()
	guestName := mux.Vars(r)["name"]
	if guestName == "" {
		ErrorResponse(w, http.StatusBadRequest, "unable to retrieve guest name from URL")
		return
	}

	var reservation model.Reservation
	if err := db.Preload("Table").Where("guest = ?", guestName).First(&reservation).Error; err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("unable to find reservation: %v", err))
		return
	}

	tmpl, err := template.New("invitation").Parse(templates.InvitationTemplate)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to load template: %v", err))
		return
	}

	buf := new(bytes.Buffer)
	tmpl.Execute(buf, struct {
		GuestName          string
		TableNumber        int
		AccompanyingGuests int
	}{
		GuestName:          reservation.Guest,
		TableNumber:        reservation.Table.Number,
		AccompanyingGuests: reservation.AccompanyingGuests,
	})

	http.StatusText(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
