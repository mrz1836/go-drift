package drift

import (
	"context"
	"fmt"
	"net/http"
)

// AttachmentData represents the raw attachment data
type AttachmentData struct {
	Data     []byte
	MimeType string
}

// GetAttachment will get the raw attachment data by its ID
// specs: https://devdocs.drift.com/docs/retrieving-a-conversations-attachments
func (c *Client) GetAttachment(ctx context.Context, attachmentID uint64) (data *AttachmentData, err error) {
	var response *RequestResponse
	if response, err = c.GetAttachmentRaw(ctx, attachmentID); err != nil {
		return nil, err
	}

	data = &AttachmentData{
		Data: response.BodyContents,
	}

	return data, nil
}

// GetAttachmentRaw will fire the HTTP request to retrieve the raw attachment data
// specs: https://devdocs.drift.com/docs/retrieving-a-conversations-attachments
func (c *Client) GetAttachmentRaw(ctx context.Context, attachmentID uint64) (*RequestResponse, error) {
	if attachmentID == 0 {
		return nil, ErrMissingAttachmentID
	}

	queryURL := fmt.Sprintf("%s/attachments/%d/data", apiEndpoint, attachmentID)
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}

// GetAttachmentFromMessage extracts attachment data from a message attachment
// This is a convenience method to get the attachment data using the attachment info from a message
func (c *Client) GetAttachmentFromMessage(ctx context.Context, attachment *MessageAttachment) (*AttachmentData, error) {
	if attachment == nil {
		return nil, ErrMissingAttachmentID
	}

	data, err := c.GetAttachment(ctx, attachment.ID)
	if err != nil {
		return nil, err
	}

	data.MimeType = attachment.MimeType
	return data, nil
}

// GetAllAttachmentsFromMessage gets all attachments from a message
func (c *Client) GetAllAttachmentsFromMessage(ctx context.Context, message *MessageData) ([]*AttachmentData, error) {
	if message == nil || len(message.Attachments) == 0 {
		return nil, nil
	}

	attachments := make([]*AttachmentData, 0, len(message.Attachments))
	for _, att := range message.Attachments {
		data, err := c.GetAttachmentFromMessage(ctx, att)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, data)
	}

	return attachments, nil
}
