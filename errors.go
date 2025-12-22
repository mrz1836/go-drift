package drift

import "errors"

// HTTP response errors
var (
	// ErrResourceNotFound is returned when a resource is not found (404).
	ErrResourceNotFound = errors.New("resource not found")

	// ErrUnauthorized is returned when oauth access token is invalid or missing (401).
	ErrUnauthorized = errors.New("oauth access token possibly invalid or missing")

	// ErrMalformedRequest is returned when request data is malformed (400).
	ErrMalformedRequest = errors.New("malformed request data")

	// ErrConflict is returned when there is an issue creating or updating a record (409).
	ErrConflict = errors.New("issue with creating or updating record, possibly already exists")

	// ErrUnexpectedStatus is returned when the status code does not match the expected status.
	ErrUnexpectedStatus = errors.New("unexpected status code")
)

// Validation errors - missing required fields
var (
	// ErrMissingAccountID is returned when account ID is empty
	ErrMissingAccountID = errors.New("account id is required")

	// ErrMissingOwnerID is returned when owner ID is empty
	ErrMissingOwnerID = errors.New("owner id is required")

	// ErrMissingContactIdentifier is returned when contact id, email or external id is missing
	ErrMissingContactIdentifier = errors.New("contact id, email or external id is required")

	// ErrMissingConversationID is returned when conversation ID is empty
	ErrMissingConversationID = errors.New("conversation id is required")

	// ErrMissingAttachmentID is returned when attachment ID is empty
	ErrMissingAttachmentID = errors.New("attachment id is required")

	// ErrMissingUserID is returned when user ID is empty
	ErrMissingUserID = errors.New("user id is required")

	// ErrMissingEmail is returned when email is empty
	ErrMissingEmail = errors.New("email is required")

	// ErrMissingMessageBody is returned when message body is empty
	ErrMissingMessageBody = errors.New("message body is required")

	// ErrMissingMessageType is returned when message type is empty
	ErrMissingMessageType = errors.New("message type is required")

	// ErrMissingClientID is returned when the client ID is empty
	ErrMissingClientID = errors.New("client id is required")

	// ErrMissingClientSecret is returned when the client secret is empty
	ErrMissingClientSecret = errors.New("client secret is required")

	// ErrMissingAccessToken is returned when the access token is empty
	ErrMissingAccessToken = errors.New("access token is required")

	// ErrMissingMinStartTime is returned when min_start_time is empty
	ErrMissingMinStartTime = errors.New("min_start_time is required")

	// ErrMissingMaxStartTime is returned when max_start_time is empty
	ErrMissingMaxStartTime = errors.New("max_start_time is required")
)

// Pagination and data errors
var (
	// ErrNoNextPage is returned when there is no next page available
	ErrNoNextPage = errors.New("no next page available")

	// ErrNoMessages is returned when no messages are found
	ErrNoMessages = errors.New("no messages found")

	// ErrTooManyUserIDs is returned when more than 20 user IDs are provided
	ErrTooManyUserIDs = errors.New("maximum of 20 user IDs allowed")
)

// StandardResponse is a common response type for simple OK/result responses
// Used by delete operations, unsubscribe, uninstall, etc.
type StandardResponse struct {
	OK     bool   `json:"ok,omitempty"`
	Result string `json:"result,omitempty"`
}
