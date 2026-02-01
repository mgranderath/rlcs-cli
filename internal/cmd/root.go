package cmd

import "github.com/alecthomas/kong"

type Context struct {
	Debug bool
}

var cli struct {
	Debug bool `help:"Enable debug mode."`

	ListTournaments ListTournamentsCmd `cmd:"" name:"list-tournaments" help:"List all tournaments for RLCS."`
}

func Execute() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
