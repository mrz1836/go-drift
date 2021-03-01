package drift

// Contact is the base contact model
type Contact struct {
	Data *contactData `json:"data"`
}

// Contacts is the multiple contacts
type Contacts struct {
	Data []*contactData `json:"data"`
}

// contactData is the internal data object
type contactData struct {
	Attributes *attributes `json:"attributes"`
	CreatedAt  int64       `json:"createdAt"`
	ID         uint64      `json:"id"`
}

// ContactFields is used for creating/updating a contact (standard attributes)
type ContactFields struct {
	Attributes *StandardAttributes `json:"attributes"`
}

// StandardAttributes are used to create new contacts
type StandardAttributes struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// attributes are the base attributes for the contact
type attributes struct {
	StandardAttributes
	CalculatedVersion                    int                    `json:"_calculated_version"`
	Classification                       string                 `json:"_classification"`
	EndUserVersion                       int                    `json:"_end_user_version"`
	Events                               map[string]interface{} `json:"events"`
	ExternalID                           string                 `json:"externalId"`
	IP                                   string                 `json:"ip"`
	LastActive                           int                    `json:"last_active"`
	LastContacted                        int                    `json:"last_contacted"`
	LastContextLocation                  string                 `json:"last_context_location"`
	OriginalConversationStartedPageTitle string                 `json:"original_conversation_started_page_title"`
	OriginalConversationStartedPageURL   string                 `json:"original_conversation_started_page_url"`
	OriginalEntrancePageTitle            string                 `json:"original_entrance_page_title"`
	OriginalEntrancePageURL              string                 `json:"original_entrance_page_url"`
	OriginalIP                           string                 `json:"original_ip"`
	OriginalRefererURL                   string                 `json:"original_referer_url"`
	RecentConversationStartedPageTitle   string                 `json:"recent_conversation_started_page_title"`
	RecentConversationStartedPageURL     string                 `json:"recent_conversation_started_page_url"`
	RecentEntrancePageTitle              string                 `json:"recent_entrance_page_title"`
	RecentEntrancePageURL                string                 `json:"recent_entrance_page_url"`
	RecentMedium                         string                 `json:"recent_medium"`
	RecentRefererURL                     string                 `json:"recent_referer_url"`
	RecentSource                         string                 `json:"recent_source"`
	SocialProfiles                       map[string]interface{} `json:"social_profiles"`
	StartDate                            int                    `json:"start_date"`
}
