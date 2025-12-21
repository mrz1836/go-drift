package drift

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHTTPGetContacts for mocking requests
type mockHTTPGetContacts struct{}

// Do is a mock http request
func (m *mockHTTPGetContacts) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	// Valid response
	if req.URL.String() == apiEndpoint+"/contacts/"+testContactID {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"data":{"id":` + testContactID + `,"createdAt":1606273669631,"attributes":{"recent_entrance_page_title":"Page Title","original_conversation_started_page_title":"Page Title","original_entrance_page_url":"https://google.com","recent_conversation_started_page_title":"Another Page Title","events":{},"phone":"` + testContactPhone + `","recent_medium":"social","_end_user_version":17899,"ip":"68.100.100.100,23.23.23.23","tags":[],"last_contacted":1613855943522,"_classification":"Engaged","recent_referer_url":"t.co","recent_source":"Twitter","socialProfiles":{},"name":"` + testContactName + `","original_referer_url":"https://googe.com","_END_USER_VERSION":17899,"_calculated_version":17899,"last_context_location":"{\"city\":\"NYC\",\"region\":\"New York\",\"country\":\"US\",\"countryName\":\"United States\",\"postalCode\":\"10901\",\"latitude\":25.5397,\"longitude\":-84.5151}","recent_conversation_started_page_url":"google.com","email":"` + testContactEmail + `","start_date":1606273669631,"original_ip":"12.12.12.12","recent_entrance_page_url":"https://google.com","externalId":"123","original_conversation_started_page_url":"google.com","original_entrance_page_title":"Page Title","last_active":1614550516644}}}`))
	} else if req.URL.String() == apiEndpoint+"/contacts/"+testContactIDBadRequest {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = io.NopCloser(nil)
	} else if req.URL.String() == apiEndpoint+"/contacts/"+testContactIDUnauthorized {
		resp.StatusCode = http.StatusUnauthorized
		resp.Body = io.NopCloser(nil)
	} else if req.URL.String() == apiEndpoint+"/contacts/"+testContactIDBadJSON {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"data":{"id":` + testContactIDBadJSON + `,"createdAt":1606273669631"attributes":{"recent_entrance_page_title""Page Title""original_conversation_started_page_title""Page Title","original_entrance_page_url":"https://google.com","recent_conversation_started_page_title":"Another Page Title","events":{},"recent_medium":"social","_end_user_version":17899,"ip":"68.100.100.100,23.23.23.23","tags":[],"last_contacted":1613855943522,"_classification":"Engaged","recent_referer_url":"t.co","recent_source":"Twitter","socialProfiles":{},"name":"` + testContactName + `","original_referer_url":"https://googe.com","_END_USER_VERSION":17899,"_calculated_version":17899,"last_context_location":"{\"city\":\"NYC\",\"region\":\"New York\",\"country\":\"US\",\"countryName\":\"United States\",\"postalCode\":\"10901\",\"latitude\":25.5397,\"longitude\":-84.5151}","recent_conversation_started_page_url":"google.com","email":"` + testContactEmail + `","start_date":1606273669631,"original_ip":"12.12.12.12","recent_entrance_page_url":"https://google.com","externalId":"123","original_conversation_started_page_url":"google.com","original_entrance_page_title":"Page Title","last_active":1614550516644}}}`))
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetContacts tests the method GetContacts()
func TestClient_GetContacts(t *testing.T) {
	t.Parallel()

	t.Run("get a valid contact by id", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPGetContacts{})

		// Create a req
		contacts, err := client.GetContacts(context.Background(), &ContactQuery{
			ID: testContactID,
		})
		require.NoError(t, err)
		assert.NotNil(t, contacts)
		assert.Len(t, contacts.Data, 1)

		// Check returned values
		assert.Equal(t, uint64(123456789), contacts.Data[0].ID)
		assert.Equal(t, int64(1606273669631), contacts.Data[0].CreatedAt)
		assert.Equal(t, testContactName, contacts.Data[0].Attributes.Name)
		assert.Equal(t, testContactEmail, contacts.Data[0].Attributes.Email)
		assert.Equal(t, testContactPhone, contacts.Data[0].Attributes.Phone)
		assert.Equal(t, "123", contacts.Data[0].Attributes.ExternalID)
		assert.Equal(t, "68.100.100.100,23.23.23.23", contacts.Data[0].Attributes.IP)
		assert.Equal(t, "12.12.12.12", contacts.Data[0].Attributes.OriginalIP)
		assert.Equal(t, "Engaged", contacts.Data[0].Attributes.Classification)
		assert.Equal(t, 1613855943522, contacts.Data[0].Attributes.LastContacted)
		assert.Equal(t, 17899, contacts.Data[0].Attributes.EndUserVersion)
		assert.Equal(t, 1614550516644, contacts.Data[0].Attributes.LastActive)
		assert.Equal(t, "social", contacts.Data[0].Attributes.RecentMedium)
		assert.Equal(t, 1606273669631, contacts.Data[0].Attributes.StartDate)
	})

	t.Run("bad request response", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPGetContacts{})

		// Create a req
		contact, err := client.GetContacts(context.Background(), &ContactQuery{
			ID: testContactIDBadRequest,
		})
		require.Error(t, err)
		assert.Nil(t, contact)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPGetContacts{})

		// Create a req
		contact, err := client.GetContacts(context.Background(), &ContactQuery{
			ID: testContactIDUnauthorized,
		})
		require.Error(t, err)
		assert.Nil(t, contact)
	})

	t.Run("bad json response", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPGetContacts{})

		// Create a req
		contact, err := client.GetContacts(context.Background(), &ContactQuery{
			ID: testContactIDBadJSON,
		})
		require.Error(t, err)
		assert.Nil(t, contact)
	})
}

// TestClient_GetContactsRaw tests the method GetContactsRaw()
func TestClient_GetContactsRaw(t *testing.T) {
	t.Parallel()

	t.Run("invalid query", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPGetContacts{})

		// Create a req
		response, err := client.GetContactsRaw(context.Background(), &ContactQuery{})
		assert.Nil(t, response)
		assert.Error(t, err)
	})

	t.Run("get a valid contact by id", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPGetContacts{})

		// Create a req
		response, err := client.GetContactsRaw(context.Background(), &ContactQuery{
			ID: testContactID,
		})
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		// Check returned values
		assert.Equal(t, apiEndpoint+"/contacts/"+testContactID, response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Len(t, response.BodyContents, 1172)
	})
}

// TestContactQuery_HasMultipleResults tests the method HasMultipleResults()
func TestContactQuery_HasMultipleResults(t *testing.T) {
	t.Parallel()

	t.Run("query with single result", func(t *testing.T) {
		q := &ContactQuery{
			ID: testContactID,
		}
		assert.False(t, q.HasMultipleResults())
	})

	t.Run("query with multiple results (email)", func(t *testing.T) {
		q := &ContactQuery{
			Email: testContactEmail,
		}
		assert.True(t, q.HasMultipleResults())
	})

	t.Run("query with multiple results (external id)", func(t *testing.T) {
		q := &ContactQuery{
			ExternalID: testContactEmail,
		}
		assert.True(t, q.HasMultipleResults())
	})
}

// TestContactQuery_BuildURL tests the method BuildURL()
func TestContactQuery_BuildURL(t *testing.T) {
	t.Parallel()

	t.Run("requires an identifier to search", func(t *testing.T) {
		q := &ContactQuery{}
		queryURL, err := q.BuildURL()
		require.Error(t, err)
		assert.Empty(t, queryURL)
	})

	t.Run("sets a limit to 1 if not given", func(t *testing.T) {
		q := &ContactQuery{ID: testContactID}
		queryURL, err := q.BuildURL()
		require.NoError(t, err)
		assert.Equal(t, 1, q.Limit)
		assert.Equal(t, apiEndpoint+"/contacts/"+testContactID, queryURL)
	})

	t.Run("url by contact id", func(t *testing.T) {
		q := &ContactQuery{ID: testContactID}
		queryURL, err := q.BuildURL()
		require.NoError(t, err)
		assert.Equal(t, apiEndpoint+"/contacts/"+testContactID, queryURL)
	})

	t.Run("url by contact email", func(t *testing.T) {
		q := &ContactQuery{Email: testContactEmail}
		queryURL, err := q.BuildURL()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(apiEndpoint+"/contacts?email="+testContactEmail+"&limit=%d", q.Limit), queryURL)
	})

	t.Run("url by contact external id", func(t *testing.T) {
		q := &ContactQuery{ExternalID: testContactEmail}
		queryURL, err := q.BuildURL()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(apiEndpoint+"/contacts?idType=external&id="+testContactEmail+"&limit=%d", q.Limit), queryURL)
	})

	t.Run("custom limit", func(t *testing.T) {
		q := &ContactQuery{Email: testContactEmail, Limit: 123}
		queryURL, err := q.BuildURL()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(apiEndpoint+"/contacts?email="+testContactEmail+"&limit=%d", 123), queryURL)
	})
}

// BenchmarkClient_GetContacts benchmarks the GetContacts method
func BenchmarkClient_GetContacts(b *testing.B) {
	client := newTestClient(&mockHTTPCreateContact{})
	fields := &ContactQuery{
		ID: testContactID,
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.GetContacts(context.Background(), fields)
	}
}

// BenchmarkClient_GetContacts benchmarks the GetContactsRaw method
func BenchmarkClient_GetContactsRaw(b *testing.B) {
	client := newTestClient(&mockHTTPCreateContact{})
	fields := &ContactQuery{
		ID: testContactID,
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.GetContactsRaw(context.Background(), fields)
	}
}
