package main

import "time"

// region api-def

type UberAPI interface {
	DriverByID(string) (*Driver, error)
	RiderByID(string) (*Rider, error)
	TripByID(string) (*Trip, error)
}

type Driver struct {
	ID   string
	Name string
}

type Rider struct {
	ID   string
	Name string
}

type Trip struct {
	ID       string
	DriverID string
	RiderID  string
}

// endregion api-def

// region impl
type fakeUberClient struct{}

func (*fakeUberClient) TripByID(id string) (*Trip, error) {
	time.Sleep(200 * time.Millisecond)
	return &Trip{
		ID:       id,
		DriverID: "42",
		RiderID:  "57",
	}, nil
}

func (*fakeUberClient) DriverByID(id string) (*Driver, error) {
	time.Sleep(500 * time.Millisecond)
	return &Driver{
		ID:   id,
		Name: "Eleanor Nelson",
	}, nil
}

func (*fakeUberClient) RiderByID(id string) (*Rider, error) {
	time.Sleep(300 * time.Millisecond)
	return &Rider{
		ID:   id,
		Name: "Richard Dickson",
	}, nil
}

// endregion impl
