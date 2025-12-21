package drift

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ContactQuery is how we want to get a contact(s)
type ContactQuery struct {
	Email      string `json:"email"`
	ExternalID string `json:"external_id"`
	ID         string `json:"id"`
	Limit      int    `json:"limit"`
}

// BuildURL will build a url depending on our query params
func (q *ContactQuery) BuildURL() (queryURL string, err error) {
	// Make sure we have something to search for
	if len(q.ID) == 0 && len(q.Email) == 0 && len(q.ExternalID) == 0 {
		err = fmt.Errorf("contact id, email or external id is required")
		return queryURL, err
	}

	// Set a default limit if no limit is given
	if q.Limit == 0 {
		q.Limit = 1
	}

	// Got an ID (highest priority)
	if len(q.ID) > 0 {
		queryURL = apiEndpoint + "/contacts/" + q.ID
	} else if len(q.Email) > 0 { // Next is email
		queryURL = fmt.Sprintf("%s/contacts?email=%s&limit=%d", apiEndpoint, q.Email, q.Limit)
	} else if len(q.ExternalID) > 0 { // Next is external id
		queryURL = fmt.Sprintf("%s/contacts?idType=external&id=%s&limit=%d", apiEndpoint, q.ExternalID, q.Limit)
	}
	return queryURL, err
}

// HasMultipleResults will return true if the query will produce multiple contacts
func (q *ContactQuery) HasMultipleResults() bool {
	return len(q.ID) == 0
}

// GetContacts will get the contact data, but then parse into a standard contact (no custom attributes)
// specs: https://devdocs.drift.com/docs/retrieving-contact
func (c *Client) GetContacts(ctx context.Context, query *ContactQuery) (contacts *Contacts, err error) {
	// Create and fire the request
	var response *RequestResponse
	if response, err = c.GetContactsRaw(
		ctx, query,
	); err != nil {
		return contacts, err
	}

	// Determine if single or multiple
	contacts = new(Contacts)
	if query.HasMultipleResults() {
		if err = json.Unmarshal(
			response.BodyContents, &contacts,
		); err != nil {
			contacts = nil
			return contacts, err
		}
	} else { // Parse as a single contact
		contact := new(Contact)
		if err = json.Unmarshal(
			response.BodyContents, &contact,
		); err != nil {
			contacts = nil
			return contacts, err
		}
		contacts.Data = append(contacts.Data, contact.Data)
	}

	return contacts, err
}

// GetContactsRaw will fire the HTTP request to retrieve the raw contact data
// specs: https://devdocs.drift.com/docs/retrieving-contact
func (c *Client) GetContactsRaw(ctx context.Context, query *ContactQuery) (response *RequestResponse, err error) {
	var queryURL string
	if queryURL, err = query.BuildURL(); err != nil {
		return response, err
	}
	if response = httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	); response.Error != nil {
		err = response.Error
	}
	return response, err
}
