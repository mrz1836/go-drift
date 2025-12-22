package drift

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateContact will fire the HTTP request to create a new contact
// specs: https://devdocs.drift.com/docs/creating-a-contact
func (c *Client) CreateContact(ctx context.Context, attributes *ContactFields) (contact *Contact, err error) {
	// Create and fire the request
	var response *RequestResponse
	if response, err = c.CreateContactRaw(
		ctx, attributes,
	); err != nil {
		return contact, err
	}

	// Parse the request
	err = response.UnmarshalTo(&contact)
	return contact, err
}

// CreateContactRaw will create a contact using custom attributes
// specs: https://devdocs.drift.com/docs/creating-a-contact
func (c *Client) CreateContactRaw(ctx context.Context, attributes interface{}) (*RequestResponse, error) {
	return c.createOrUpdateContact(ctx, 0, attributes)
}

// createOrUpdateContact will create or update a contact
func (c *Client) createOrUpdateContact(ctx context.Context, contactID uint64,
	attributes interface{},
) (response *RequestResponse, err error) {
	// Marshall the attributes
	var data []byte
	if data, err = json.Marshal(attributes); err != nil {
		return response, err
	}

	// Set the method based on the type of request
	method := http.MethodPost
	endpointURL := apiEndpoint + "/contacts"
	if contactID > 0 { // Update if contact id is passed
		method = http.MethodPatch
		endpointURL = fmt.Sprintf(apiEndpoint+"/contacts/%d", contactID)
	}

	// Create and fire the request
	if response = httpRequest(
		ctx, c, &httpPayload{
			Data:           data,
			ExpectedStatus: http.StatusOK,
			Method:         method,
			URL:            endpointURL,
		},
	); response.Error != nil {
		err = response.Error
	}
	return response, err
}
