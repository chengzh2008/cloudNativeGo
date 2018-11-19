package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../lib/persistence"
	"github.com/gorilla/mux"
)

type eventServicesHandler struct {
	dbhandler persistence.DatabaseHandler
}

func NewEventHandler(databasehandler persistence.DatabaseHandler) *eventServicesHandler {
	return &eventServicesHandler{
		dbhandler: databasehandler,
	}
}

func (esh *eventServicesHandler) FindEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["searchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search criteria found, you can either search by id via /id/4 or search by name via /name/<yourname>"}`)
		return
	}
	searchValue, ok := vars["searchValue"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search value found, you can either search by id via /id/4 or search by name via /name/<yourname>"}`)
		return
	}

	var event persistence.Event
	var err error

	switch strings.ToLower(criteria) {
	case "name":
		event, err = esh.dbhandler.FindEventByName(searchValue)
	case "id":
		id, err := hex.DecodeString(searchValue)
		if err == nil {
			event, err = esh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		fmt.Fprintf(w, `{"error": "%s"}`, err)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)

}

func (esh *eventServicesHandler) AllEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := esh.dbhandler.FindAllAvailableEvents()
	fmt.Println(events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "Error occured while trying to find all available events %s"}`, err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&events)
}

func (esh *eventServicesHandler) NewEventHandler(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "Error occured while decoding event data %s"}`, err)
		return
	}
	id, err := esh.dbhandler.AddEvent(event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "Error occured while persisting event data %d  %s"}`, id, err)
		return
	}
	fmt.Fprint(w, `{"id": %d}`, id)
}
