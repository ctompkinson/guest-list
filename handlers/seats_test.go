package handlers

import (
	"github.com/gorilla/mux"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandleGetEmptySeats(t *testing.T) {
	database.Init()
	db := database.Get()

	database.ClearAndCreate()

	tables := []model.Table{
		{Number: 1, Seats: 10},
		{Number: 1, Seats: 10},
		{Number: 1, Seats: 10},
	}

	now := time.Now()
	reservations := []model.Reservation{
		{Guest: "Bob", AccompanyingGuests: 5, Table: tables[0], ArrivalTime: &now},
		{Guest: "Taylor", AccompanyingGuests: 2, Table: tables[0], ArrivalTime: &now},
		{Guest: "Scott", AccompanyingGuests: 9, Table: tables[2], ArrivalTime: &now},
		{Guest: "Yokie", AccompanyingGuests: 9, Table: tables[1]}, // This one hasnt arrived so the seats are still free
	}

	for _, reservation := range reservations {
		db.Create(&reservation)
	}

	req, err := http.NewRequest("GET", "/seats_empty", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/seats_empty", HandleGetEmptySeats)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"seats_empty":21}`, strings.TrimSpace(rr.Body.String()))
}
