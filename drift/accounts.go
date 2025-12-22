package drift

// Account is the single account response wrapper
type Account struct {
	Data *accountData `json:"data"`
}

// Accounts is the list accounts response wrapper
type Accounts struct {
	Data *accountsData `json:"data"`
}

// accountsData is the internal list response structure
type accountsData struct {
	Accounts []*accountData `json:"accounts"`
	Total    int            `json:"total"`
	Next     string         `json:"next,omitempty"`
}

// accountData is the internal account data object
type accountData struct {
	AccountID        string            `json:"accountId"`
	OwnerID          uint64            `json:"ownerId"`
	Name             string            `json:"name,omitempty"`
	Domain           string            `json:"domain,omitempty"`
	Deleted          bool              `json:"deleted"`
	Targeted         bool              `json:"targeted"`
	CreateDateTime   int64             `json:"createDateTime"`
	UpdateDateTime   int64             `json:"updateDateTime"`
	CustomProperties []*CustomProperty `json:"customProperties,omitempty"`
}

// CustomProperty represents a custom property on an account
type CustomProperty struct {
	Label string      `json:"label"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"` // STRING, EMAIL, NUMBER, PHONE, URL, DATE, DATETIME, ENUM, ENUMARRAY, LATLON, LAT, LON, TEAMMEMBER
}

// AccountFields is used for creating/updating an account
type AccountFields struct {
	AccountID        string            `json:"accountId,omitempty"` // Required for update
	OwnerID          uint64            `json:"ownerId"`             // Required
	Name             string            `json:"name,omitempty"`
	Domain           string            `json:"domain,omitempty"`
	Targeted         bool              `json:"targeted,omitempty"`
	CustomProperties []*CustomProperty `json:"customProperties,omitempty"`
}

// AccountListQuery contains pagination parameters for listing accounts
type AccountListQuery struct {
	Index int // Starting index (default: 0)
	Size  int // Batch size (default: 10, max: 65)
}
