package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ctompkinson/guest-list/database"
	"github.com/ctompkinson/guest-list/handlers"
	"net/http"
	"time"
)

type Server struct {
	router *mux.Router
}

func Start() {
	router := mux.NewRouter()

	router.HandleFunc("/table/{tableNumber}", handlers.HandleCreateTable).Methods("POST")
	router.HandleFunc("/table/{tableNumber}", handlers.HandleDeleteTable).Methods("DELETE")
	router.HandleFunc("/table/{tableNumber}", handlers.HandleGetTable).Methods("GET")

	router.HandleFunc("/guest_list/{name}", handlers.HandleCreateReservation).Methods("POST")
	router.HandleFunc("/guest_list/{name}", handlers.HandleDeleteReservation).Methods("DELETE")
	router.HandleFunc("/guest_list", handlers.HandleGetReservations).Methods("GET")

	router.HandleFunc("/guests", handlers.HandleListGuests).Methods("GET")
	router.HandleFunc("/guest/{name}", handlers.HandleGuestArrival).Methods("PUT")
	// We reuse delete reservation because its effectively the same thing
	router.HandleFunc("/guest/{name}", handlers.HandleDeleteReservation).Methods("DELETE")

	router.HandleFunc("/seats_empty", handlers.HandleGetEmptySeats).Methods("GET")
	router.HandleFunc("/invitation/{name}", handlers.HandleCreateInvitation).Methods("GET")

	// Start the database and and make it available to the API
	if err := database.Init(); err != nil {
		panic(err)
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8080", // TODO: Make the port and address configurable
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting Server")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("error starting server:", err)
	}
}
