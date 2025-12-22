package drift

// GDPRRequest is the request payload for GDPR operations.
type GDPRRequest struct {
	Email string `json:"email"`
}

// GDPRRetrievalResponse is the response from the GDPR retrieval endpoint.
// The retrieval is processed asynchronously and results are emailed to the org owner.
type GDPRRetrievalResponse struct {
	Data *gdprRetrievalData `json:"data"`
}

// gdprRetrievalData contains the retrieval response data.
type gdprRetrievalData struct {
	Message     string `json:"message"`
	SentToEmail string `json:"sentToEmail"`
}

// GDPRDeletionResponse is the response from the GDPR deletion endpoint.
// The deletion is processed asynchronously and is permanent.
type GDPRDeletionResponse struct {
	Data *gdprDeletionData `json:"data"`
}

// gdprDeletionData contains the deletion response data.
type gdprDeletionData struct {
	Message string `json:"message"`
}
