package listener

import (
	"log"
	"fmt"
	"../../contracts"
	"../../lib/msgqueue"
	"../../lib/persistence"
	"gopkg.in/mgo.v2/bson"
)

type EventProcessor struct {
	EventListener msgqueue.EventListener
	Database persistence.DatabaseHandler
}

func (p *EventProcessor) ProcessEvents() error {
	log.Println("Listening to events")
	log.Println(p.EventListener)

	received, errors, err := p.EventListener.Listen("event.created")
	if err != nil {
		log.Printf("error while start listening: %s", err)
		return err
	}

	for {
		select {
		case evt := <-received:
			fmt.Printf("got event %T: %s\n", evt, evt)
			p.handleEvent(evt)
		case err = <-errors:
			log.Printf("received error while processing msg: %s", err)
		}
	}
}

func (p *EventProcessor) handleEvent(event msgqueue.Event) {
	log.Printf("handling event: ")
	switch e := event.(type) {
	case *contracts.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)
		p.Database.AddEvent(persistence.Event{ID: bson.ObjectId(e.ID), Name: e.Name})
	case *contracts.LocationCreatedEvent:
		log.Printf("location %s created: %s", e.ID, e)
		p.Database.AddLocation(persistence.Location{ID: bson.ObjectId(e.ID)})
	default:
		log.Printf("unknown event: %t", e)
	}
}
