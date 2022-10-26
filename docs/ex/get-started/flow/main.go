// region directive
//go:build cff

package main

// endregion directive

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.uber.org/cff"
)

// region fake-client
var uber UberAPI = new(fakeUberClient)

// endregion fake-client

// region resp-decl
type Response struct {
	Rider    string
	Driver   string
	HomeCity string
}

// endregion resp-decl

// region main
func main() {
	// endregion main
	// region resp-var
	var res *Response
	// endregion resp-var
	// region ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// endregion ctx
	// region flow-start
	err := cff.Flow(ctx,
		// endregion flow-start
		// region flow-dots
		// ...
		// endregion flow-dots
		// region params
		// region resp-var
		cff.Params(12),
		// endregion params
		cff.Results(&res),
		// region get-trip
		cff.Task(func(tripID int) (*Trip, error) {
			// endregion resp-var
			return uber.TripByID(tripID)
		}),
		// endregion get-trip
		// region get-driver
		cff.Task(func(trip *Trip) (*Driver, error) {
			return uber.DriverByID(trip.DriverID)
		}),
		// endregion get-driver
		// region get-rider
		cff.Task(func(trip *Trip) (*Rider, error) {
			return uber.RiderByID(trip.RiderID)
		}),
		// endregion get-rider
		// region get-location
		cff.Task(func(rider *Rider) (*Location, error) {
			return uber.LocationByID(rider.HomeID)
		}),
		// endregion get-location
		// region last-task
		// region error
		cff.Task(func(r *Rider, d *Driver, home *Location) *Response {
			return &Response{
				Driver:   d.Name,
				Rider:    r.Name,
				HomeCity: home.City,
			}
		}),
		// endregion last-task
	)
	// region tail
	if err != nil {
		log.Fatal(err)
	}
	// endregion error

	fmt.Println(res.Driver, "drove", res.Rider, "who lives in", res.HomeCity)
	// endregion tail
}
