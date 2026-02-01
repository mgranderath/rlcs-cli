package cmd

import "github.com/alecthomas/kong"

type Context struct {
	Debug bool
}

var cli struct {
	Debug   bool             `help:"Enable debug mode."`
	Version kong.VersionFlag `name:"version" short:"v" help:"Show version and exit."`

	ListTournaments ListTournamentsCmd `cmd:"" name:"list-tournaments" help:"List all tournaments for RLCS."`
	GetBrackets     GetBracketsCmd     `cmd:"" name:"get-brackets" help:"Get tournament brackets for a specific tournament."`
	GetMatches      GetMatchesCmd      `cmd:"" name:"get-matches" help:"Get all matches for a specific tournament."`
}

func Execute(version string) {
	ctx := kong.Parse(&cli, kong.Vars{
		"version": version,
	})
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
