package rest

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"../lib/persistence"
)

func ServeAPI(endpoint string, databasehandler persistence.DatabaseHandler) error {
	handler := NewEventHandler(databasehandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{searchCriteria}/{searchValue}").HandlerFunc(handler.FindEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)
	fmt.Println("starting serving...")
	return http.ListenAndServe(endpoint, r)
}
