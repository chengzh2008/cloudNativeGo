package kafka

import (
	"encoding/json"
	"../../msgqueue"
	"github.com/Shopify/sarama"
)

type kafkaEventEmitter struct {
	producer sarama.SyncProducer
}

func NewKafkaEventEmitter(client sarama.Client) (msgqueue.EventEmitter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := &kafkaEventEmitter{
		producer: producer,
	}
	return emitter, nil
}

func (k *kafkaEventEmitter) Emit(event msgqueue.Event) error {
	envelope := messageEnvelope{
		event.EventName(),
		event,
	}
	jsonDoc, err := json.Marshal(&envelope)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: event.EventName(),
		Value: sarama.ByteEncoder(jsonDoc),
	}
	_, _, err = k.producer.SendMessage(msg)
	return err
}
