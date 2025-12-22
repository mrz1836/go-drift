package drift

// Conversation is the base conversation model
type Conversation struct {
	Data *conversationData `json:"data"`
}

// Conversations is for multiple conversations (list response)
type Conversations struct {
	Data  []*conversationData `json:"data"`
	Links *PaginationLinks    `json:"links,omitempty"`
}

// conversationData is the internal data object
type conversationData struct {
	ContactID         uint64             `json:"contactId"`
	ConversationTags  []*ConversationTag `json:"conversationTags,omitempty"`
	CreatedAt         int64              `json:"createdAt"`
	ID                uint64             `json:"id"`
	InboxID           int                `json:"inboxId"`
	OrgID             int                `json:"orgId,omitempty"`
	Participants      []uint64           `json:"participants"`
	RelatedPlaybookID int                `json:"relatedPlaybookId,omitempty"`
	Status            string             `json:"status"`
	UpdatedAt         int64              `json:"updatedAt"`
}

// ConversationTag represents a tag on a conversation
type ConversationTag struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

// PaginationLinks for paginated responses
type PaginationLinks struct {
	Next string `json:"next,omitempty"`
	Self string `json:"self,omitempty"`
}

// ConversationStats is the response for bulk conversation stats
type ConversationStats struct {
	ConversationCount map[string]int `json:"conversationCount"`
}

// NewConversationRequest for creating a new conversation
type NewConversationRequest struct {
	Email   string                  `json:"email"`
	Message *NewConversationMessage `json:"message"`
}

// NewConversationMessage is the initial message for a new conversation
type NewConversationMessage struct {
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	Body       string                 `json:"body"`
}

// TranscriptResponse for the transcript endpoint
type TranscriptResponse struct {
	Data string `json:"data"`
}
