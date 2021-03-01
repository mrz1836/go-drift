package drift

import (
	"context"
	"encoding/json"
)

// UpdateContact will fire the HTTP request to update an existing contact
// specs: https://devdocs.drift.com/docs/creating-a-contact
func (c *Client) UpdateContact(ctx context.Context, contactID uint64,
	attributes *ContactFields) (contact *Contact, err error) {

	// Create and fire the request
	var response *RequestResponse
	if response, err = c.UpdateContactRaw(
		ctx, contactID, attributes,
	); err != nil {
		return
	}

	// Parse the request
	err = json.Unmarshal(response.BodyContents, &contact)
	return

}

// UpdateContactRaw will update an existing contact using a custom attribute struct
// specs: https://devdocs.drift.com/docs/updating-a-contact
func (c *Client) UpdateContactRaw(ctx context.Context, contactID uint64,
	attributes interface{}) (*RequestResponse, error) {
	return c.createOrUpdateContact(ctx, contactID, attributes)
}
