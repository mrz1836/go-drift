package drift

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHTTPUpdateContact for mocking requests
type mockHTTPUpdateContact struct{}

// Do is a mock http request
func (m *mockHTTPUpdateContact) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if req.URL.String() == apiEndpoint+"/contacts/"+testContactID {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"data":{"id":` + testContactID + `,"createdAt":1606273669631,"attributes":{"recent_entrance_page_title":"Page Title","original_conversation_started_page_title":"Page Title","original_entrance_page_url":"https://google.com","recent_conversation_started_page_title":"Another Page Title","events":{},"phone":"` + testContactPhone + `","recent_medium":"social","_end_user_version":17899,"ip":"68.100.100.100,23.23.23.23","tags":[],"last_contacted":1613855943522,"_classification":"Engaged","recent_referer_url":"t.co","recent_source":"Twitter","socialProfiles":{},"name":"` + testContactName + `2","original_referer_url":"https://googe.com","_END_USER_VERSION":17899,"_calculated_version":17899,"last_context_location":"{\"city\":\"NYC\",\"region\":\"New York\",\"country\":\"US\",\"countryName\":\"United States\",\"postalCode\":\"10901\",\"latitude\":25.5397,\"longitude\":-84.5151}","recent_conversation_started_page_url":"google.com","email":"` + testContactEmail + `","start_date":1606273669631,"original_ip":"12.12.12.12","recent_entrance_page_url":"https://google.com","externalId":"123","original_conversation_started_page_url":"google.com","original_entrance_page_title":"Page Title","last_active":1614550516644}}}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_UpdateContact tests the method UpdateContact()
func TestClient_UpdateContact(t *testing.T) {
	t.Parallel()

	t.Run("update a standard contact", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPUpdateContact{})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		assert.NoError(t, err)

		// Create a req
		var contact *Contact
		contact, err = client.UpdateContact(
			context.Background(), id,
			&ContactFields{&StandardAttributes{
				Name: testContactName + "2",
			}})
		assert.NotNil(t, contact)
		assert.NoError(t, err)

		// Got a contact
		assert.Equal(t, id, contact.Data.ID)
		assert.Equal(t, int64(1606273669631), contact.Data.CreatedAt)
		assert.Equal(t, testContactName+"2", contact.Data.Attributes.Name)
	})
}

// BenchmarkClient_UpdateContact benchmarks the UpdateContact method
func BenchmarkClient_UpdateContact(b *testing.B) {
	client := newTestClient(&mockHTTPCreateContact{})
	id, _ := strconv.ParseUint(testContactID, 10, 64)
	fields := &ContactFields{&StandardAttributes{
		Email: testContactEmail,
		Name:  testContactName,
		Phone: testContactPhone,
	}}
	for i := 0; i < b.N; i++ {
		_, _ = client.UpdateContact(context.Background(), id, fields)
	}
}
