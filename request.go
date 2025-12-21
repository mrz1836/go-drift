package drift

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// RequestResponse is the response from a request
type RequestResponse struct {
	BodyContents []byte `json:"body_contents"` // Raw body response
	Error        error  `json:"error"`         // If an error occurs
	Method       string `json:"method"`        // Method is the HTTP method used
	PostData     string `json:"post_data"`     // PostData is the post data submitted if POST/PUT request
	StatusCode   int    `json:"status_code"`   // StatusCode is the last code from the request
	URL          string `json:"url"`           // URL is used for the request
}

// httpPayload is used for a httpRequest
type httpPayload struct {
	Data           []byte `json:"data"`
	ExpectedStatus int    `json:"expected_status"`
	Method         string `json:"method"`
	URL            string `json:"url"`
}

// httpRequest is a generic request wrapper that can be used without constraints
func httpRequest(ctx context.Context, client *Client,
	payload *httpPayload,
) (response *RequestResponse) {
	// Set reader
	var bodyReader io.Reader

	// Start the response
	response = new(RequestResponse)

	// Add post data if applicable
	if payload.Method == http.MethodPost || payload.Method == http.MethodPatch {
		bodyReader = bytes.NewBuffer(payload.Data)
		response.PostData = string(payload.Data)
	}

	// Store for debugging purposes
	response.Method = payload.Method
	response.URL = payload.URL

	// Start the request
	var request *http.Request
	if request, response.Error = http.NewRequestWithContext(
		ctx, payload.Method, payload.URL, bodyReader,
	); response.Error != nil {
		return response
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", client.Options.UserAgent)

	// Set the content type on Method
	if payload.Method == http.MethodPost || payload.Method == http.MethodPatch {
		request.Header.Set("Content-Type", "application/json")
	}

	// Set an access token if supplied
	if len(client.OAuthAccessToken) > 0 {
		request.Header.Set("Authorization", "Bearer "+client.OAuthAccessToken)
	}

	// Fire the http request
	var resp *http.Response
	if resp, response.Error = client.httpClient.Do(request); response.Error != nil {
		if resp != nil {
			response.StatusCode = resp.StatusCode
		}
		return response
	}

	// Close the response body
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	// Set the status
	response.StatusCode = resp.StatusCode

	// Check status code
	if payload.ExpectedStatus != resp.StatusCode {
		switch resp.StatusCode {
		case http.StatusNotFound:
			response.Error = fmt.Errorf("resource not found: %s", response.URL)
		case http.StatusUnauthorized:
			response.Error = fmt.Errorf("oauth access token possible invalid or missing")
		case http.StatusBadRequest:
			response.Error = fmt.Errorf("malformatted request data")
		case http.StatusConflict:
			response.Error = fmt.Errorf("issue with creating or updating record, possibly already exists")
		default:
			response.Error = fmt.Errorf(
				"status code: %d does not match %d",
				resp.StatusCode, payload.ExpectedStatus,
			)
		}
		return response
	}

	// Read the body
	response.BodyContents, response.Error = ioutil.ReadAll(resp.Body)

	return response
}
