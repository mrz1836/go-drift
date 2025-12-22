package drift

// User is the base user model (single user response)
type User struct {
	Data *userData `json:"data"`
}

// Users is the multiple users response (list endpoint)
type Users struct {
	Data []*userData `json:"data"`
}

// UsersMap is the multiple users response (get by IDs endpoint - map structure)
type UsersMap struct {
	Data map[string]*userData `json:"data"`
}

// userData is the internal user data object
type userData struct {
	ID           uint64 `json:"id"`
	OrgID        uint64 `json:"orgId"`
	Name         string `json:"name"`
	Alias        string `json:"alias"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Locale       string `json:"locale"`
	Availability string `json:"availability"` // AVAILABLE, OFFLINE, ON_CALL
	Role         string `json:"role"`         // member, admin, agent
	TimeZone     string `json:"timeZone"`
	AvatarURL    string `json:"avatarUrl"`
	Verified     bool   `json:"verified"`
	Bot          bool   `json:"bot"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
}

// UserUpdateFields is used for updating a user (PATCH request)
type UserUpdateFields struct {
	Name         string `json:"name,omitempty"`
	Alias        string `json:"alias,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Locale       string `json:"locale,omitempty"`
	AvatarURL    string `json:"avatarUrl,omitempty"`
	Availability string `json:"availability,omitempty"` // AVAILABLE or OFFLINE
}

// Meeting is the meeting model for booked meetings
type Meeting struct {
	AgentID         uint64 `json:"agentId"`
	OrgID           uint64 `json:"orgId"`
	Status          string `json:"status"`
	MeetingSource   string `json:"meetingSource"`
	SchedulerID     int64  `json:"schedulerId"`
	EventID         string `json:"eventId"`
	Slug            string `json:"slug"`
	SlotStart       int64  `json:"slotStart"`
	SlotEnd         int64  `json:"slotEnd"`
	UpdatedAt       int64  `json:"updatedAt"`
	ScheduledAt     int64  `json:"scheduledAt"`
	MeetingType     string `json:"meetingType"`
	ConversationID  int64  `json:"conversationId"`
	EndUserTimeZone string `json:"endUserTimeZone"`
	MeetingNotes    string `json:"meetingNotes"`
	BookedBy        uint64 `json:"bookedBy"`
	ConferenceType  string `json:"conferenceType"`
	IsRecurring     bool   `json:"isRecurring"`
	IsPrivate       bool   `json:"isPrivate"`
}

// Meetings is the response for booked meetings
type Meetings struct {
	Data []*Meeting `json:"data"`
}
