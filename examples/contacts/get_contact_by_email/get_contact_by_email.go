// Package main provides an example of retrieving a contact by email from the Drift API.
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

	// Get a "standard" contact by id (no custom attributes)
	contacts, err := client.GetContacts(
		context.Background(), &drift.ContactQuery{
			Email: os.Getenv("TEST_DRIFT_CONTACT_EMAIL"),
		},
	)
	if err != nil {
		log.Fatal("failed: ", err.Error()) //nolint:gosec // G706: example code, values from trusted API response
		return
	}

	// See the standard contact data
	log.Println(contacts.Data[0].ID)               //nolint:gosec // G706: example code, values from trusted API response
	log.Println(contacts.Data[0].CreatedAt)        //nolint:gosec // G706: example code, values from trusted API response
	log.Println(contacts.Data[0].Attributes)       //nolint:gosec // G706: example code, values from trusted API response
	log.Println(contacts.Data[0].Attributes.Email) //nolint:gosec // G706: example code, values from trusted API response
}
