package drift

// Team is the base team model (single team response)
type Team struct {
	Data *teamData `json:"data"`
}

// Teams is the multiple teams response (list endpoint)
type Teams struct {
	Data []*teamData `json:"data"`
}

// teamData is the internal team data object
type teamData struct {
	ID                   uint64   `json:"id"`
	OrgID                uint64   `json:"orgId"`
	WorkspaceID          string   `json:"workspaceId"`
	Name                 string   `json:"name"`
	UpdatedAt            int64    `json:"updatedAt"`
	Members              []uint64 `json:"members"`
	Owner                uint64   `json:"owner"`
	Status               string   `json:"status"` // ENABLED, ARCHIVED
	Main                 bool     `json:"main"`
	AutoOffline          bool     `json:"autoOffline"`
	TeamCsatEnabled      bool     `json:"teamCsatEnabled"`
	TeamAvailabilityMode string   `json:"teamAvailabilityMode"` // ALWAYS_ONLINE, ALWAYS_OFFLINE, CUSTOM_HOURS
	ResponseTimerEnabled bool     `json:"responseTimerEnabled"`
}
