package commands_structs

import (
	"AVMRuntime/src/aasm/commands_handlers/assembler"
	"fmt"
)

type AssembleCommand struct {
	Input  string `arg:"" type:"path" help:"Input file path which should be a valid AASM assembly file (UTF-8)."`
	Output string `short:"o" type:"path" help:"Defines the output path of the assembled file. Defaults to the current directory, and the name of the input file '.avme' extension."`
}

func (cmd *AssembleCommand) Run() error {
	fmt.Println("Assembling file: ", cmd.Input)

	fileData, fileDataError := assembler.ReadFile(cmd.Input)

	if fileDataError != nil {
		return fmt.Errorf("couldn't read input file '%s': %s", cmd.Input, fileDataError)
	}

	fileDatString := string(fileData)

	firstLine, firstLineError := assembler.ExtractAndCheckFirstLine(fileDatString)

	if firstLineError != nil {
		return firstLineError
	}

	fmt.Println(firstLine)

	return nil
}
