// Package main provides an example of updating a contact in the Drift API.
package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/mrz1836/go-drift/drift"
)

func main() {
	// Create a new client
	client := drift.NewClient(
		os.Getenv("TEST_DRIFT_OAUTH_TOKEN"), nil, nil,
	)

	// Parse our env string into a number (just for this example)
	id, _ := strconv.ParseUint(os.Getenv("TEST_DRIFT_CONTACT_ID"), 10, 64)

	// Update a standard contact
	contact, err := client.UpdateContact(
		context.Background(), id, &drift.ContactFields{
			Attributes: &drift.StandardAttributes{
				Name: "John Doe",
			},
		},
	)
	if err != nil {
		log.Fatal("failed: ", err.Error()) //nolint:gosec // G706: example code, values from trusted API response
		return
	}

	// See the standard contact data
	log.Println(contact.Data.ID)         //nolint:gosec // G706: example code, values from trusted API response
	log.Println(contact.Data.CreatedAt)  //nolint:gosec // G706: example code, values from trusted API response
	log.Println(contact.Data.Attributes) //nolint:gosec // G706: example code, values from trusted API response
}
