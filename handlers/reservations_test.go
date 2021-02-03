package handlers

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleCreateReservation(t *testing.T) {
	database.Init()
	db := database.Get()

	cases := []struct {
		name             string
		url              string
		body             string
		expectedStatus   int
		expectedResponse string
		createTable      *model.Table
	}{
		{
			"good",
			"/guest_list/bob",
			`{ "table": 1, "accompanying_guests": 5 }`,
			http.StatusOK,
			`{"name":"bob"}`,
			&model.Table{Number: 1, Seats: 6},
		},
		{
			"noData",
			"/guest_list/bob",
			``,
			http.StatusBadRequest,
			`{"message":"unable to parse body: EOF"}`,
			nil,
		},
		{
			"noTable",
			"/guest_list/bob",
			`{ "table": 1, "accompanying_guests": 5 }`,
			http.StatusInternalServerError,
			`{"message":"failed to find table: record not found"}`,
			nil,
		},
		{
			"noSeats",
			"/guest_list/bob",
			`{ "table": 1, "accompanying_guests": 5 }`,
			http.StatusInternalServerError,
			`{"message":"not enough seats available on selected table"}`,
			&model.Table{Number: 1, Seats: 5},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database.ClearAndCreate()

			if c.createTable != nil {
				db.Create(c.createTable)
			}

			req, err := http.NewRequest("POST", c.url, bytes.NewBuffer([]byte(c.body)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/guest_list/{name}", HandleCreateReservation)
			router.ServeHTTP(rr, req)

			assert.Equal(t, c.expectedStatus, rr.Code)
			assert.Equal(t, c.expectedResponse, strings.TrimSpace(rr.Body.String()))
			router.ServeHTTP(rr, req)
		})
	}
}

func TestHandleCreateReservation_Duplicate(t *testing.T) {
	database.Init()
	database.ClearAndCreate()
	db := database.Get()

	db.Create(&model.Table{Number: 1, Seats: 10})

	req, err := http.NewRequest("POST", "/guest_list/bob",
		bytes.NewBuffer([]byte(`{ "table": 1, "accompanying_guests": 5 }`)))
	require.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/guest_list/{name}", HandleCreateReservation)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req)
	assert.Equal(t, http.StatusBadRequest, rr2.Code)
}

func TestHandleDeleteReservation(t *testing.T) {
	database.Init()
	db := database.Get()

	cases := []struct {
		name              string
		url               string
		expectedStatus    int
		createReservation *model.Reservation
		createTable       *model.Table
	}{
		{
			"good",
			"/guest_list/bob",
			http.StatusOK,
			&model.Reservation{Guest: "bob", AccompanyingGuests: 1},
			&model.Table{Number: 1, Seats: 5},
		},
		{
			"noReservation",
			"/guest_list/bob",
			http.StatusBadRequest,
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
			if c.createReservation != nil {
				c.createReservation.Table = *c.createTable
				db.Create(c.createReservation)
			}

			req, err := http.NewRequest("DELETE", c.url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/guest_list/{name}", HandleDeleteReservation)
			router.ServeHTTP(rr, req)

			assert.Equal(t, c.expectedStatus, rr.Code)
			router.ServeHTTP(rr, req)
		})
	}
}

func TestHandleGetReservations(t *testing.T) {
	database.Init()
	db := database.Get()

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
			"/guest_list",
			http.StatusOK,
			`{"guests":[{"name":"bob","table":1,"accompanying_guests":1},{"name":"taylor","table":1,"accompanying_guests":2}]}`,
			[]*model.Reservation{
				{Guest: "bob", AccompanyingGuests: 1},
				{Guest: "taylor", AccompanyingGuests: 2},
			},
			&model.Table{Number: 1, Seats: 5},
		},
		{
			"noReservation",
			"/guest_list",
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
			router.HandleFunc("/guest_list", HandleGetReservations)
			router.ServeHTTP(rr, req)

			assert.Equal(t, c.expectedStatus, rr.Code)
			assert.Equal(t, c.expectedResponse, strings.TrimSpace(rr.Body.String()))
			router.ServeHTTP(rr, req)
		})
	}
}
