package domain

import "time"

// Bracket represents a tournament bracket
type Bracket struct {
	TournamentUUID         string
	TournamentName         string
	ParentTournamentName   string
	ParentTournamentFormat string
	CircuitName            string
	StartDate              time.Time
	EndDate                time.Time
	Index                  int
	Label                  string
	Format                 string
	NumberOfTeams          *int
	Matches                []Match
}

// Match represents a match in a bracket
type Match struct {
	UUID         string
	Type         string
	Index        int
	Name         string
	TimeOfSeries time.Time
	TeamA        MatchTeam
	TeamB        MatchTeam
	TeamAScore   int
	TeamBScore   int
	Maps         []MatchMap
	ExternalID   string
	WinnerGoesTo *BracketDestination
	LoserGoesTo  *BracketDestination
	IsLive       bool
	IsCompleted  bool
}

// MatchTeam represents a team in a match
type MatchTeam struct {
	UUID         string
	Name         string
	Shorthand    string
	Location     string
	IsEliminated bool
}

// MatchMap represents a single map/game in a match
type MatchMap struct {
	UUID               string
	ScheduledStartTime time.Time
	ActualStartTime    time.Time
	Name               string
	MatchEndedTime     time.Time
	TeamAScore         int
	TeamBScore         int
	ExternalID         string
}

// BracketDestination represents where a team advances
type BracketDestination struct {
	TournamentUUID  string
	SeriesUUID      string
	BracketPosition string
}
