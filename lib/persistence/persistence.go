package persistence

// DatabaseHandler interface: decoupling database handler implementation
type DatabaseHandler interface {
	AddEvent(Event) ([]byte, error)
	AddLocation(Location) (Location, error)
	AddBookingForUser([]byte, Booking) error
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllAvailableEvents() ([]Event, error)
}
