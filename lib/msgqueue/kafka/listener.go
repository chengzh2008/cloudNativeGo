package kafka

import (
	"../../msgqueue"
	"log"
	"fmt"
	"encoding/json"
	"github.com/Shopify/sarama"
)

type kafkaEventListener struct {
	consumer sarama.Consumer
	partitions []int32
	mapper msgqueue.EventMapper
}

func NewKafkaEventListener(client sarama.Client, partitions []int32) (msgqueue.EventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &kafkaEventListener{
		consumer: consumer,
		partitions: partitions,
	}
	return listener, nil
}

func (k *kafkaEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	var err error
	topic := "events"
	results := make(chan msgqueue.Event)
	errors := make(chan error)

	partitions := k.partitions
	if len(partitions) == 0 {
		partitions, err = k.consumer.Partitions(topic)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Printf("topic %s has partitions: %v", topic, partitions)

	for _, partition := range partitions {
		con, err := k.consumer.ConsumePartition(topic, partition, 0)
		if err != nil {
			return nil, nil, err
		}

		go func() {
			for msg := range con.Messages() {
				body := messageEnvelope{}
				err := json.Unmarshal(msg.Value, &body)
				if err != nil {
					errors <- fmt.Errorf("could not JSON-decode message: %s", err)
					continue
				}

                event, err := k.mapper.MapEvent(body.EventName, body.Payload)
                if err != nil {
                    errors <- fmt.Errorf("could not map message: %v", err)
                    continue
                }

				results <- event
			}
		}()
	}

	return results, errors, nil
}
