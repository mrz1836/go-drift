// Package main provides an example of creating a timeline event in the Drift API.
package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/mrz1836/go-drift"
)

func main() {
	// Create a new client
	client := drift.NewClient(
		os.Getenv("TEST_DRIFT_OAUTH_TOKEN"), nil, nil,
	)

	// Parse our env string into a number (just for this example)
	id, _ := strconv.ParseUint(os.Getenv("TEST_DRIFT_CONTACT_ID"), 10, 64)

	// Create a new timeline event
	event, err := client.CreateTimelineEvent(
		context.Background(), &drift.TimelineEvent{
			ContactID: id,
			Event:     "test-event-name-goes-here",
		},
	)
	if err != nil {
		log.Fatal("failed: ", err.Error())
		return
	}

	// See the standard contact data
	log.Println(event.Data.Event)
	log.Println(event.Data.CreatedAt)
	log.Println(event.Data.ContactID)
}
