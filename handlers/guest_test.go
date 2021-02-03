package handlers

import (
	"bytes"
	"fmt"
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

func TestHandleListGuests(t *testing.T) {
	database.Init()
	db := database.Get()

	now := time.Now()
	formattedTime := now.Format("02/01/06 15:04")
	fullOut := fmt.Sprintf(`{"guests":[{"name":"bob","accompanying_guests":1,"time_arrived":"%s"},{"name":"taylor","accompanying_guests":2,"time_arrived":"%s"}]}`, formattedTime, formattedTime)
	cases := []struct {
		name               string
		url                string
		expectedStatus     int
		expectedResponse   string
		createReservations []*model.Reservation
		createTable        *model.Table
	}{
		{
			"good",
			"/guests",
			http.StatusOK,
			fullOut,
			[]*model.Reservation{
				{Guest: "bob", AccompanyingGuests: 1, ArrivalTime: &now},
				{Guest: "taylor", AccompanyingGuests: 2, ArrivalTime: &now},
				{Guest: "scott", AccompanyingGuests: 2},
			},
			&model.Table{Number: 1, Seats: 5},
		},
		{
			"noReservation",
			"/guests",
			http.StatusOK,
			`{"guests":[]}`,
			nil,
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database.ClearAndCreate()

			if c.createTable != nil {
				db.Create(c.createTable)
			}
			if c.createReservations != nil {
				for _, res := range c.createReservations {
					res.Table = *c.createTable
					db.Create(res)
				}
			}

			req, err := http.NewRequest("GET", c.url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/guests", HandleListGuests)
			router.ServeHTTP(rr, req)

			assert.Equal(t, c.expectedStatus, rr.Code)
			assert.Equal(t, c.expectedResponse, strings.TrimSpace(rr.Body.String()))
			router.ServeHTTP(rr, req)
		})
	}
}

func TestHandleGuestArrival(t *testing.T) {
	database.Init()
	db := database.Get()

	cases := []struct {
		name              string
		url               string
		body              string
		expectedStatus    int
		expectedResponse  string
		createTable       *model.Table
		createReservation *model.Reservation
	}{
		{
			"good",
			"/guest/bob",
			`{ "accompanying_guests": 1 }`,
			http.StatusOK,
			`{"name":"bob"}`,
			&model.Table{Number: 1, Seats: 6},
			&model.Reservation{Guest: "bob", AccompanyingGuests: 1},
		},
		{
			"noData",
			"/guest/bob",
			``,
			http.StatusBadRequest,
			`{"message":"unable to parse body: EOF"}`,
			nil,
			nil,
		},
		{
			"noSeats",
			"/guest/bob",
			`{ "accompanying_guests": 5 }`,
			http.StatusInternalServerError,
			`{"message":"not enough seats available on selected table"}`,
			&model.Table{Number: 1, Seats: 5},
			&model.Reservation{Guest: "bob", AccompanyingGuests: 4},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database.ClearAndCreate()

			if c.createTable != nil {
				db.Create(c.createTable)
			}
			if c.createReservation != nil {
				c.createReservation.Table = *c.createTable
				db.Create(c.createReservation)
			}

			req, err := http.NewRequest("PUT", c.url, bytes.NewBuffer([]byte(c.body)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/guest/{name}", HandleGuestArrival)
			router.ServeHTTP(rr, req)

			assert.Equal(t, c.expectedStatus, rr.Code)
			assert.Equal(t, c.expectedResponse, strings.TrimSpace(rr.Body.String()))
			router.ServeHTTP(rr, req)
		})
	}
}
