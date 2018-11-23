package contracts

type EventBookedEvent struct {
	EventID string `json:"eventId"`
	UserID string `json:"userId"`
}

func (c *EventBookedEvent) EventName() string {
	return "event.booked"
}

func (c *EventBookedEvent) PartitionKey() string {
	return c.EventID
}
