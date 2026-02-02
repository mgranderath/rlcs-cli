package cmd

import "github.com/alecthomas/kong"

type Context struct {
	Debug bool
}

// TournamentsCmd groups all tournament-related commands
type TournamentsCmd struct {
	List     ListTournamentsCmd     `cmd:"" name:"list" help:"List all tournaments."`
	Matches  TournamentsMatchesCmd  `cmd:"" name:"matches" help:"List matches across tournaments."`
	Brackets TournamentsBracketsCmd `cmd:"" name:"brackets" help:"Get brackets for a specific tournament."`
}

// MatchesCmd groups all match-related commands
type MatchesCmd struct {
	List MatchesListCmd `cmd:"" name:"list" help:"List matches for a specific tournament."`
	Get  MatchesGetCmd  `cmd:"" name:"get" help:"Get detailed information for a specific match."`
}

var cli struct {
	Debug   bool             `help:"Enable debug mode."`
	Version kong.VersionFlag `name:"version" short:"v" help:"Show version and exit."`

	Tournaments TournamentsCmd `cmd:"" name:"tournaments" help:"Tournament-related commands."`
	Matches     MatchesCmd     `cmd:"" name:"matches" help:"Match-related commands."`
}

func Execute(version string) {
	ctx := kong.Parse(&cli, kong.Vars{
		"version": version,
	})
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
