package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// RetrieveGDPR triggers a GDPR data retrieval request for the given email.
// The request is processed asynchronously and results are emailed to the org owner.
// specs: https://devdocs.drift.com/docs/gdpr-retrieval
func (c *Client) RetrieveGDPR(ctx context.Context, email string) (*GDPRRetrievalResponse, error) {
	return c.RetrieveGDPRWithRequest(ctx, &GDPRRequest{Email: email})
}

// RetrieveGDPRWithRequest triggers a GDPR data retrieval using a request struct.
// specs: https://devdocs.drift.com/docs/gdpr-retrieval
func (c *Client) RetrieveGDPRWithRequest(ctx context.Context, request *GDPRRequest) (response *GDPRRetrievalResponse, err error) {
	var reqResponse *RequestResponse
	if reqResponse, err = c.RetrieveGDPRRaw(ctx, request); err != nil {
		return nil, err
	}

	err = reqResponse.UnmarshalTo(&response)
	return response, err
}

// RetrieveGDPRRaw triggers a GDPR data retrieval and returns the raw response.
// specs: https://devdocs.drift.com/docs/gdpr-retrieval
func (c *Client) RetrieveGDPRRaw(ctx context.Context, request *GDPRRequest) (*RequestResponse, error) {
	if request == nil {
		return nil, ErrMissingEmail
	}
	if err := requireString(request.Email, ErrMissingEmail); err != nil {
		return nil, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response := httpRequest(ctx, c, &httpPayload{
		Data:           data,
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodPost,
		URL:            apiEndpoint + "/gdpr/retrieve",
	})

	return response, response.Error
}
