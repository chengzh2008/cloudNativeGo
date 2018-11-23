package amqp

import (
	"../../msgqueue"
	"../../../contracts"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type amqpEventListener struct {
	connection *amqp.Connection
	exchange string
	queue string
}

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	return err
}

func NewAMQPEventListener(conn *amqp.Connection, exchange string, queue string) (msgqueue.EventListener, error) {
	listener := &amqpEventListener{
		connection: conn,
		exchange: exchange,
		queue: queue,
	}

	err := listener.setup()
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (a *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer channel.Close()
	for _, eventName := range eventNames {
		if err := channel.QueueBind(a.queue, eventName, "events", false, nil); err != nil {
			return nil, nil, err
		}
	}

	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	eventChan := make(chan msgqueue.Event)
	errorChan := make(chan error)
	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errorChan <- fmt.Errorf("msg did not contain x-event-name header")
				msg.Nack(false, false)
				continue
			}
			eventName, ok := rawEventName.(string)
			if !ok {
				errorChan <- fmt.Errorf("x-event-name header is not string, but %t", rawEventName)
				msg.Nack(false, false)
				continue
			}

			var event msgqueue.Event
			switch eventName {
			case "event.created":
				eventChan <- new(contracts.EventCreatedEvent)
			default:
				errorChan <- fmt.Errorf("event type %s is unknown", eventName)
				continue
			}

			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				errorChan <- err
				continue
			}
			eventChan <- event
			msg.Ack(false)
		}
	}()
	return eventChan, errorChan, nil
}
