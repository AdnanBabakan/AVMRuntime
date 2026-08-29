package main

import (
	"AVMRuntime/src/aasm/cli_structs"

	"github.com/alecthomas/kong"
)

func main() {
	var cli cli_structs.AASMCLI

	ctx := kong.Parse(&cli)
	err := ctx.Run()

	ctx.FatalIfErrorf(err)
}
