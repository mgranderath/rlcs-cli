// Package blast contains API models for the Blast.tv API
package blast

// Circuit represents a circuit in the Blast API
type Circuit struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	GameID string `json:"gameId"`
}

// Tournament represents a tournament response from the Blast API
type Tournament struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	StartDate     string                 `json:"startDate"`
	EndDate       string                 `json:"endDate"`
	CircuitID     string                 `json:"circuitId"`
	Circuit       Circuit                `json:"circuit"`
	PrizePool     string                 `json:"prizePool"`
	Location      string                 `json:"location"`
	NumberOfTeams int                    `json:"numberOfTeams"`
	ExternalID    *string                `json:"externalId"`
	Region        string                 `json:"region"`
	Grouping      string                 `json:"grouping"`
	Description   string                 `json:"description"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     string                 `json:"createdAt"`
	DeletedAt     *string                `json:"deletedAt"`
	UpdatedAt     string                 `json:"updatedAt"`
}
