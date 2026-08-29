package main

import (
	"AVMRuntime/src/aasm/clistructs"

	"github.com/alecthomas/kong"
)

func main() {
	var cli clistructs.AASMCLI

	ctx := kong.Parse(&cli)
	err := ctx.Run()

	ctx.FatalIfErrorf(err)
}
