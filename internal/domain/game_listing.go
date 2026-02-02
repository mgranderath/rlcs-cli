package domain

// GameListing represents a match along with its tournament context
type GameListing struct {
	TournamentID   string
	TournamentName string
	Match          Match
}
