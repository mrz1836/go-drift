// Package main provides an example of creating a contact in the Drift API.
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

	// Create a standard contact
	contact, err := client.CreateContact(
		context.Background(), &drift.ContactFields{
			Attributes: &drift.StandardAttributes{
				Email: "john@email.com",
				Name:  "John Doe",
				Phone: "15554443333",
			},
		},
	)
	if err != nil {
		log.Fatal("failed: ", err.Error())
		return
	}

	// See the standard contact data
	log.Println(contact.Data.ID)
	log.Println(contact.Data.CreatedAt)
	log.Println(contact.Data.Attributes)
}
