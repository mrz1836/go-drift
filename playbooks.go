package drift

// Playbook is the base playbook model (single playbook response)
type Playbook struct {
	Data *playbookData `json:"data"`
}

// Playbooks is the multiple playbooks response (list endpoint)
type Playbooks struct {
	Data []*playbookData `json:"data"`
}

// playbookData is the internal playbook data object
type playbookData struct {
	ID              uint64                 `json:"id"`
	Name            string                 `json:"name"`
	OrgID           uint64                 `json:"orgId"`
	Meta            map[string]interface{} `json:"meta"`
	CreatedAt       int64                  `json:"createdAt"`
	UpdatedAt       int64                  `json:"updatedAt"`
	CreatedAuthorID uint64                 `json:"createdAuthorId"`
	UpdatedAuthorID uint64                 `json:"updatedAuthorId"`
	InteractionID   uint64                 `json:"interactionId"`
	ReportType      string                 `json:"reportType"`
	Goals           []*playbookGoal        `json:"goals"`
}

// playbookGoal is a goal within a playbook
type playbookGoal struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// ConversationalLandingPage represents a conversational landing page
type ConversationalLandingPage struct {
	PlaybookID     uint64 `json:"playbookId"`
	PlaybookName   string `json:"playbookName"`
	LandingPageURL string `json:"landingPageUrl"`
}

// ConversationalLandingPages is the response for listing conversational landing pages
type ConversationalLandingPages struct {
	Data []*ConversationalLandingPage `json:"data"`
}
