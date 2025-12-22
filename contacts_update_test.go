package drift

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockUpdateContact returns a mock for contact update operations
func mockUpdateContact() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":{"id":`+testContactID+`,"createdAt":1606273669631,"attributes":{"recent_entrance_page_title":"Page Title","original_conversation_started_page_title":"Page Title","original_entrance_page_url":"https://google.com","recent_conversation_started_page_title":"Another Page Title","events":{},"phone":"`+testContactPhone+`","recent_medium":"social","_end_user_version":17899,"ip":"68.100.100.100,23.23.23.23","tags":[],"last_contacted":1613855943522,"_classification":"Engaged","recent_referer_url":"t.co","recent_source":"Twitter","socialProfiles":{},"name":"`+testContactName+`2","original_referer_url":"https://googe.com","_END_USER_VERSION":17899,"_calculated_version":17899,"last_context_location":"{\"city\":\"NYC\",\"region\":\"New York\",\"country\":\"US\",\"countryName\":\"United States\",\"postalCode\":\"10901\",\"latitude\":25.5397,\"longitude\":-84.5151}","recent_conversation_started_page_url":"google.com","email":"`+testContactEmail+`","start_date":1606273669631,"original_ip":"12.12.12.12","recent_entrance_page_url":"https://google.com","externalId":"123","original_conversation_started_page_url":"google.com","original_entrance_page_title":"Page Title","last_active":1614550516644}}}`),
	)
}

// TestClient_UpdateContact tests the method UpdateContact()
func TestClient_UpdateContact(t *testing.T) {
	t.Parallel()

	t.Run("update a standard contact", func(t *testing.T) {
		client := newTestClient(mockUpdateContact())

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		var contact *Contact
		contact, err = client.UpdateContact(
			context.Background(), id,
			&ContactFields{&StandardAttributes{
				Name: testContactName + "2",
			}})
		require.NoError(t, err)
		assert.NotNil(t, contact)

		// Got a contact
		assert.Equal(t, id, contact.Data.ID)
		assert.Equal(t, int64(1606273669631), contact.Data.CreatedAt)
		assert.Equal(t, testContactName+"2", contact.Data.Attributes.Name)
	})

	t.Run("returns error when UpdateContactRaw fails", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		contact, err := client.UpdateContact(
			context.Background(), id,
			&ContactFields{&StandardAttributes{
				Name: testContactName,
			}})

		require.Error(t, err)
		assert.Nil(t, contact)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 404 not found", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusNotFound))

		contact, err := client.UpdateContact(
			context.Background(), 999999,
			&ContactFields{&StandardAttributes{
				Name: testContactName,
			}})

		require.Error(t, err)
		assert.Nil(t, contact)
		assert.ErrorIs(t, err, ErrResourceNotFound)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":{"invalid json`))

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		contact, err := client.UpdateContact(
			context.Background(), id,
			&ContactFields{&StandardAttributes{
				Name: testContactName,
			}})

		require.Error(t, err)
		assert.Nil(t, contact)
	})
}

// TestClient_UpdateContactRaw tests the method UpdateContactRaw()
func TestClient_UpdateContactRaw(t *testing.T) {
	t.Parallel()

	t.Run("updates contact successfully", func(t *testing.T) {
		client := newTestClient(mockUpdateContact())

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.UpdateContactRaw(
			context.Background(), id,
			&ContactFields{&StandardAttributes{
				Name: testContactName + "2",
			}})

		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPatch, response.Method)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.UpdateContactRaw(
			context.Background(), id,
			&ContactFields{&StandardAttributes{
				Name: testContactName,
			}})

		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("uses correct endpoint URL with contact ID", func(t *testing.T) {
		client := newTestClient(mockUpdateContact())

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.UpdateContactRaw(
			context.Background(), id,
			&ContactFields{&StandardAttributes{
				Name: testContactName,
			}})

		require.NoError(t, err)
		assert.Contains(t, response.URL, "/contacts/"+testContactID)
	})
}

// BenchmarkClient_UpdateContact benchmarks the UpdateContact method
func BenchmarkClient_UpdateContact(b *testing.B) {
	client := newTestClient(mockCreateContact())
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
