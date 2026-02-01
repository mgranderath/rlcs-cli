package blast

// Bracket represents a tournament bracket from the Blast API
type Bracket struct {
	TournamentUUID         string      `json:"tournamentUuid"`
	TournamentName         string      `json:"tournamentName"`
	ParentTournamentName   string      `json:"parentTournamentName"`
	ParentTournamentFormat string      `json:"parentTournamentFormat"`
	CircuitName            string      `json:"circuitName"`
	StartDate              string      `json:"startDate"`
	EndDate                string      `json:"endDate"`
	Index                  int         `json:"index"`
	Label                  string      `json:"label"`
	Format                 string      `json:"format"`
	NumberOfTeams          *int        `json:"numberOfTeams"`
	Metadata               interface{} `json:"metadata"`
	Matches                []Match     `json:"matches"`
}

// Match represents a match in a tournament bracket
type Match struct {
	UUID         string              `json:"uuid"`
	Type         string              `json:"type"`
	Index        int                 `json:"index"`
	Name         string              `json:"name"`
	TimeOfSeries string              `json:"timeOfSeries"`
	TeamA        Team                `json:"teamA"`
	TeamB        Team                `json:"teamB"`
	TeamAScore   int                 `json:"teamAScore"`
	TeamBScore   int                 `json:"teamBScore"`
	Maps         []Map               `json:"maps"`
	Metadata     MatchMetadata       `json:"metadata"`
	ExternalID   string              `json:"externalId"`
	WinnerGoesTo *BracketDestination `json:"winnerGoesTo"`
	LoserGoesTo  *BracketDestination `json:"loserGoesTo"`
	IsLive       bool                `json:"isLive"`
	IsCompleted  bool                `json:"isCompleted"`
}

// Team represents a team in a match
type Team struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Shorthand    string `json:"shorthand"`
	Location     string `json:"location"`
	IsEliminated bool   `json:"isEliminated"`
}

// Map represents a single game/map in a match
type Map struct {
	UUID               string `json:"uuid"`
	ScheduledStartTime string `json:"scheduledStartTime"`
	ActualStartTime    string `json:"actualStartTime"`
	Name               string `json:"name"`
	MatchEndedTime     string `json:"matchEndedTime"`
	TeamAScore         int    `json:"teamAScore"`
	TeamBScore         int    `json:"teamBScore"`
	ExternalID         string `json:"externalId"`
}

// MatchMetadata contains additional match information
type MatchMetadata struct {
	Type              string `json:"_t"`
	TeamBlueTeamID    string `json:"teamBlueTeamId"`
	TeamOrangeTeamID  string `json:"teamOrangeTeamId"`
	ExternalStreamURL string `json:"externalStreamUrl"`
}

// BracketDestination represents where a team advances after winning/losing
type BracketDestination struct {
	TournamentUUID  string `json:"tournamentUUID"`
	SeriesUUID      string `json:"seriesUUID"`
	BracketPosition string `json:"bracketPosition"`
}
