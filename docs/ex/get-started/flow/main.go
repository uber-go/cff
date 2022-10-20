//go:build cff

package main

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
	Rider  string
	Driver string
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
		// region params
		// region resp-var
		cff.Params("1234"),
		// endregion params
		cff.Results(&res),
		// region get-trip
		cff.Task(func(tripID string) (*Trip, error) {
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
		// region last-task
		// ...
		// region tail
		cff.Task(func(r *Rider, d *Driver) *Response {
			return &Response{
				Driver: d.Name,
				Rider:  r.Name,
			}
		}),
		// endregion last-task
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Driver, "drove", res.Rider)
	// endregion tail
}
