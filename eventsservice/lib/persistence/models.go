package persistence

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

// User model
type User struct {
	ID       bson.ObjectId `bson:"_id"`
	First    string
	Last     string
	Age      int
	Bookings []Booking
}

func (u *User) String() string {
	return fmt.Sprintf("id: %s, first_name: %s, last_name: %s, Age: %d, Bookings: %v", u.ID, u.First, u.Last, u.Age, u.Bookings)
}

// Booking model
type Booking struct {
	Date    int64
	EventID []byte
	Seats   int
}

// Event model
type Event struct {
	ID bson.ObjectId `bson:"_id"`
	Name string
	Duration int
	StartDate int64
	EndDate int64
	Location Location
}

// Location model
type Location struct {
	ID bson.ObjectId `bson:"_id"`
	Name string
	Address string
	Country string
	OpenTime int
	CloseTime int
	Halls []Hall
}

// Hall model
type Hall struct {
	Name string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int `json:"capacity"`
}
