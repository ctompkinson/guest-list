package model

import (
	"gorm.io/gorm"
	"time"
)

type Reservation struct {
	gorm.Model         `json:"-"`
	Guest              string `gorm:"unique"`
	AccompanyingGuests int
	TableID            int
	Table              Table
	ArrivalTime        *time.Time
}

type FormattedReservation struct {
	Guest              string `json:"name"`
	Table              int    `json:"table"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

type FormattedGuestArrival struct {
	Guest              string `json:"name"`
	AccompanyingGuests int    `json:"accompanying_guests"`
	TimeArrived        string `json:"time_arrived"`
}

// FormatAsReservation creates a simple string representation of a reservation without arrival time as only a checked
// in guest has one
func (r *Reservation) FormatAsReservation() FormattedReservation {
	return FormattedReservation{
		Guest:              r.Guest,
		Table:              r.Table.Number,
		AccompanyingGuests: r.AccompanyingGuests,
	}
}

// FormatAsGuestArrival creates a simple string representation of a reservation with arrival time and without table
// to match API spec
func (r *Reservation) FormatAsGuestArrival() FormattedGuestArrival {
	return FormattedGuestArrival{
		Guest:              r.Guest,
		AccompanyingGuests: r.AccompanyingGuests,
		TimeArrived:        r.ArrivalTime.Format("02/01/06 15:04"),
	}
}
