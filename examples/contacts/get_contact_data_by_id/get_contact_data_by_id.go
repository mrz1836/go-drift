package main

import (
	"context"
	"log"
	"os"

	"github.com/mrz1836/go-drift"
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
		log.Fatal("failed: ", data.Error.Error())
		return
	}

	// See the raw contact data
	log.Println(string(data.BodyContents))
}
