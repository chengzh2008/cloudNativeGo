package rest

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"../../lib/persistence"
)

func ServeAPI(endpoint string, tlsendpoint string, databasehandler persistence.DatabaseHandler) (chan error, chan error) {
	handler := NewEventHandler(databasehandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{searchCriteria}/{searchValue}").HandlerFunc(handler.FindEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)
	fmt.Println("starting serving...")

	httpErrChan := make(chan error)
	httpTLSErrChan := make(chan error)
	go func() {
		httpErrChan <- http.ListenAndServe(endpoint, r)
	}()
	go func() {
		httpTLSErrChan <- http.ListenAndServeTLS(tlsendpoint, "./eventsservice/rest/cert.pem", "./eventsservice/rest/key.pem", r)
	}()
	return httpErrChan, httpTLSErrChan
}