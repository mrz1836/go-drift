package drift

// Message is the base message model
type Message struct {
	Data *MessageData `json:"data"`
}

// Messages is for multiple messages response
type Messages struct {
	Data       *MessagesListData   `json:"data"`
	Pagination *MessagesPagination `json:"pagination,omitempty"`
}

// MessagesListData wraps the messages array
type MessagesListData struct {
	Messages []*MessageData `json:"messages"`
}

// MessagesPagination contains pagination info for messages
type MessagesPagination struct {
	Next string `json:"next,omitempty"`
}

// MessageData is the data object for a message
type MessageData struct {
	Attachments    []*MessageAttachment   `json:"attachments,omitempty"`
	Attributes     map[string]interface{} `json:"attributes,omitempty"`
	Author         *MessageAuthor         `json:"author"`
	Body           string                 `json:"body,omitempty"`
	Buttons        []*MessageButton       `json:"buttons,omitempty"`
	Context        *MessageContext        `json:"context,omitempty"`
	ConversationID uint64                 `json:"conversationId"`
	CreatedAt      int64                  `json:"createdAt"`
	ID             uint64                 `json:"id"`
	OrgID          int                    `json:"orgId"`
	Type           string                 `json:"type"` // "chat" or "private_note"
}

// MessageAuthor represents who sent the message
type MessageAuthor struct {
	Bot  bool   `json:"bot,omitempty"`
	ID   uint64 `json:"id"`
	Type string `json:"type"` // "contact" or "user"
}

// MessageButton for interactive buttons in messages
type MessageButton struct {
	Label    string          `json:"label"`
	Reaction *ButtonReaction `json:"reaction,omitempty"`
	Style    string          `json:"style,omitempty"`
	Type     string          `json:"type"`
	Value    string          `json:"value"`
}

// ButtonReaction defines what happens when a button is clicked
type ButtonReaction struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// MessageContext contains metadata about the message
type MessageContext struct {
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
}

// MessageAttachment represents a file attachment in a message
type MessageAttachment struct {
	FileName string `json:"fileName"`
	ID       uint64 `json:"id"`
	MimeType string `json:"mimeType"`
	URL      string `json:"url"`
}

// CreateMessageRequest for creating a new message in a conversation
type CreateMessageRequest struct {
	Body    string           `json:"body,omitempty"`
	Buttons []*MessageButton `json:"buttons,omitempty"`
	Type    string           `json:"type"` // "chat" or "private_note"
	UserID  uint64           `json:"userId,omitempty"`
}
