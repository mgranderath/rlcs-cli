package main

import (
	"github.com/mgranderath/rlcs-cli/internal/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
