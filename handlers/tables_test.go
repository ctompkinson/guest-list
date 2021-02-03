package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleCreateTable(t *testing.T) {
	database.Init()

	cases := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{"good", "/table/1", http.StatusOK},
		{"junk", "/table/junk", http.StatusBadRequest},
		{"missing", "/table", http.StatusNotFound},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database.ClearAndCreate()

			body, err := json.Marshal(createTableRequest{Seats: 10})
			require.NoError(t, err)
			req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/table/{tableNumber}", HandleCreateTable)
			router.ServeHTTP(rr, req)

			assert.Equal(t, c.expectedStatus, rr.Code)
			router.ServeHTTP(rr, req)
		})
	}
}

func Test_handleCreateTable_Duplicate(t *testing.T) {
	database.Init()
	database.ClearAndCreate()

	body, err := json.Marshal(createTableRequest{Seats: 10})
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/table/1", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/table/{tableNumber}", HandleCreateTable)
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	router.ServeHTTP(rr, req)

	rr2 := httptest.NewRecorder()
	router2 := mux.NewRouter()
	router2.HandleFunc("/table/{tableNumber}", HandleCreateTable)
	router2.ServeHTTP(rr2, req)
	assert.Equal(t, http.StatusInternalServerError, rr2.Code)
	router.ServeHTTP(rr2, req)
}

func TestHandleDeleteTable(t *testing.T) {
	database.Init()
	db := database.Get()

	cases := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{"good", "/table/1", http.StatusOK},
		{"junk", "/table/junk", http.StatusInternalServerError},
		{"missing", "/table", http.StatusNotFound},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database.ClearAndCreate()
			db.Create(&model.Table{Number: 1, Seats: 0})

			req, err := http.NewRequest("DELETE", c.url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/table/{tableNumber}", HandleDeleteTable)
			router.ServeHTTP(rr, req)

			assert.Equal(t, c.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetTable(t *testing.T) {
	database.Init()
	db := database.Get()

	cases := []struct {
		name           string
		url            string
		expectedStatus int
		expectBody     bool
	}{
		{"good", "/table/1", http.StatusOK, true},
		{"junk", "/table/junk", http.StatusInternalServerError, false},
		{"missing", "/table", http.StatusNotFound, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database.ClearAndCreate()
			db.Create(&model.Table{Number: 1, Seats: 10})

			req, err := http.NewRequest("GET", c.url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/table/{tableNumber}", HandleGetTable)
			router.ServeHTTP(rr, req)

			require.Equal(t, c.expectedStatus, rr.Code)
			if c.expectBody {
				var table model.Table
				b := rr.Body.String()
				fmt.Println(b)
				err := json.Unmarshal(rr.Body.Bytes(), &table)
				assert.NoError(t, err)
				assert.Equal(t, 1, table.Number)
				assert.Equal(t, 10, table.Seats)
			}
		})
	}
}
