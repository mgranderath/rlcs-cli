package blast

// MatchResponse represents a match from the matches endpoint
type MatchResponse struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	ScheduledAt string                `json:"scheduledAt"`
	Type        string                `json:"type"`
	Index       int                   `json:"index"`
	ExternalID  string                `json:"externalId"`
	Circuit     Circuit               `json:"circuit"`
	Tournament  Tournament            `json:"tournament"`
	Stage       Stage                 `json:"stage"`
	TeamA       MatchResponseTeam     `json:"teamA"`
	TeamB       MatchResponseTeam     `json:"teamB"`
	TeamAScore  int                   `json:"teamAScore"`
	TeamBScore  int                   `json:"teamBScore"`
	Maps        []MatchResponseMap    `json:"maps"`
	Metadata    MatchResponseMetadata `json:"metadata"`
}

// MatchResponseTeam represents a team in a match from the matches endpoint
type MatchResponseTeam struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	ShortName   string      `json:"shortName"`
	Nationality string      `json:"nationality"`
	ExternalID  *string     `json:"externalId"`
	Metadata    interface{} `json:"metadata"`
}

// MatchResponseMap represents a single game/map in a match from the matches endpoint
type MatchResponseMap struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ScheduledAt string `json:"scheduledAt"`
	StartedAt   string `json:"startedAt"`
	EndedAt     string `json:"endedAt"`
	ExternalID  string `json:"externalId"`
	TeamAScore  int    `json:"teamAScore"`
	TeamBScore  int    `json:"teamBScore"`
}

// MatchResponseMetadata represents metadata in a match from the matches endpoint
type MatchResponseMetadata struct {
	T                 string `json:"_t"`
	TeamBlueTeamID    string `json:"teamBlueTeamId"`
	TeamOrangeTeamID  string `json:"teamOrangeTeamId"`
	ExternalStreamURL string `json:"externalStreamUrl"`
}

// Stage represents a tournament stage
type Stage struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Format        string      `json:"format"`
	NumberOfTeams *int        `json:"numberOfTeams"`
	Metadata      interface{} `json:"metadata"`
	StartDate     string      `json:"startDate"`
	EndDate       string      `json:"endDate"`
	Index         int         `json:"index"`
}
