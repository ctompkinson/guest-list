package handlers

import (
	"github.com/gorilla/mux"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleCreateInvitation(t *testing.T) {
	database.Init()
	db := database.Get()

	database.ClearAndCreate()

	tables := []model.Table{
		{Number: 1, Seats: 10},
	}

	now := time.Now()
	reservations := []model.Reservation{
		{Guest: "Bob", AccompanyingGuests: 5, Table: tables[0], ArrivalTime: &now},
	}

	for _, reservation := range reservations {
		db.Create(&reservation)
	}

	req, err := http.NewRequest("GET", "/invitation/bob", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/invitation/{name}", HandleCreateInvitation)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// TODO: Should check contents properly instead of just accepting not empty
	assert.NotEmpty(t, rr.Body.String())
}
