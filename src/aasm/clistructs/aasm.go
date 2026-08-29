package clistructs

import "AVMRuntime/src/aasm/clistructs/commands_structs"

type AASMCLI struct {
	Assemble commands_structs.AssembleCommand `cmd:"" aliases:"asm,a" help:"Assembles an AASM aasmfile resulting in an executable binary."`
}
