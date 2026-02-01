package domain

import "time"

// Region represents the geographical region of a tournament
type Region string

const (
	RegionNA   Region = "NA"
	RegionEU   Region = "EU"
	RegionAPAC Region = "APAC"
	RegionSAM  Region = "SAM"
	RegionOCE  Region = "OCE"
	RegionMENA Region = "MENA"
	RegionSSA  Region = "SSA"
	RegionNone Region = "" // For majors and world championships
)

// TournamentType categorizes tournaments by their level
type TournamentType string

const (
	TypeOpen              TournamentType = "Open"
	TypeMajor             TournamentType = "Major"
	TypeWorldChampionship TournamentType = "WorldChampionship"
	TypeKickoff           TournamentType = "Kickoff"
)

// Tournament is the domain model for RLCS tournaments
type Tournament struct {
	ID          string
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	CircuitID   string
	PrizePool   string
	Location    string
	TeamCount   int
	Region      Region
	Type        TournamentType
	Description string
	IsOnline    bool
	IsMajor     bool
}

// IsUpcoming returns true if the tournament hasn't started yet
func (t *Tournament) IsUpcoming(now time.Time) bool {
	return t.StartDate.After(now)
}

// IsOngoing returns true if the tournament is currently running
func (t *Tournament) IsOngoing(now time.Time) bool {
	return !t.StartDate.After(now) && !t.EndDate.Before(now)
}

// IsPast returns true if the tournament has ended
func (t *Tournament) IsPast(now time.Time) bool {
	return t.EndDate.Before(now)
}

// Duration returns the duration of the tournament
func (t *Tournament) Duration() time.Duration {
	return t.EndDate.Sub(t.StartDate)
}
