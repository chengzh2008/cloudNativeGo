package rest

import (
	"net/http"
	"time"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"../../lib/msgqueue"
	"../../lib/persistence"
)

func ServeAPI(listenAddr string, database persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) {
	r := mux.NewRouter()
	r.Methods("POST").Path("/events/{eventID}/bookings").Handler(&CreateBookingHandler{eventEmitter, database})
	fmt.Println("starting serving...")

	server := handlers.CORS()(r)
	srv := http.Server{
		Handler: server,
		Addr: listenAddr,
		WriteTimeout: 2 * time.Second,
		ReadTimeout: 1 * time.Second,
	}

	srv.ListenAndServe()
}
