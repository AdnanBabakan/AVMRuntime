package cli_structs

import "AVMRuntime/src/aasm/cli_structs/commands_structs"

type AASMCLI struct {
	Assemble commands_structs.AssembleCommand `cmd:"" aliases:"asm,a" help:"Assembles an AASM file resulting in an executable binary."`
}
