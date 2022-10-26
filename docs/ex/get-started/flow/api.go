package main

import "time"

// region api-def

type UberAPI interface {
	DriverByID(int) (*Driver, error)
	RiderByID(int) (*Rider, error)
	TripByID(int) (*Trip, error)
	LocationByID(int) (*Location, error)
}

type Driver struct {
	ID   int
	Name string
}

type Location struct {
	ID    int
	City  string
	State string
	// ...
}

type Rider struct {
	ID     int
	Name   string
	HomeID int
}

type Trip struct {
	ID       int
	DriverID int
	RiderID  int
}

// endregion api-def

var _ UberAPI = (*fakeUberClient)(nil)

// region impl
type fakeUberClient struct{}

func (*fakeUberClient) DriverByID(id int) (*Driver, error) {
	time.Sleep(500 * time.Millisecond)
	return &Driver{
		ID:   id,
		Name: "Eleanor Nelson",
	}, nil
}

func (*fakeUberClient) LocationByID(id int) (*Location, error) {
	time.Sleep(200 * time.Millisecond)
	return &Location{
		ID:    id,
		City:  "San Francisco",
		State: "California",
	}, nil
}

func (*fakeUberClient) RiderByID(id int) (*Rider, error) {
	time.Sleep(300 * time.Millisecond)
	return &Rider{
		ID:   id,
		Name: "Richard Dickson",
	}, nil
}

func (*fakeUberClient) TripByID(id int) (*Trip, error) {
	time.Sleep(150 * time.Millisecond)
	return &Trip{
		ID:       id,
		DriverID: 42,
		RiderID:  57,
	}, nil
}

// endregion impl
