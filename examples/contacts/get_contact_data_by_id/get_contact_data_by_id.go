// Package main provides an example of retrieving raw contact data by ID from the Drift API.
package main

import (
	"context"
	"log"
	"os"

	"github.com/mrz1836/go-drift/drift"
)

func main() {
	// Create a new client
	client := drift.NewClient(
		os.Getenv("TEST_DRIFT_OAUTH_TOKEN"), nil, nil,
	)

	// Get a contact by id (raw data)
	data, err := client.GetContactsRaw(
		context.Background(), &drift.ContactQuery{
			ID: os.Getenv("TEST_DRIFT_CONTACT_ID"),
		},
	)
	if err != nil {
		log.Fatal("failed: ", data.Error.Error()) //nolint:gosec // G706: example code, values from trusted API response
		return
	}

	// See the raw contact data
	log.Println(string(data.BodyContents)) //nolint:gosec // G706: example code, values from trusted API response
}
