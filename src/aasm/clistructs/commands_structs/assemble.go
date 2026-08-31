package commands_structs

import (
	"AVMRuntime/src/aasm/commandshandlers/assembler/aasmfile"
	"fmt"
)

type AssembleCommand struct {
	Input  string `arg:"" type:"path" help:"Input aasmfile path which should be a valid AASM assembly aasmfile (UTF-8)."`
	Output string `short:"o" type:"path" help:"Defines the output path of the assembled aasmfile. Defaults to the current directory, and the name of the input aasmfile '.avme' extension."`
}

func (cmd *AssembleCommand) Run() error {
	fmt.Println("Assembling aasmfile: ", cmd.Input)

	fmt.Println("Checking the first line...")

	fileData, fileDataError := aasmfile.ReadFile(cmd.Input)

	if fileDataError != nil {
		return fmt.Errorf("couldn't read input aasmfile '%s': %s", cmd.Input, fileDataError)
	}

	fileDatString := string(fileData)

	firstLine, firstLineError := aasmfile.ExtractAndCheckFirstLine(fileDatString)

	if firstLineError != nil {
		return firstLineError
	}

	fmt.Printf("Detected AASM file version: %s\n", firstLine.VersionString())

	fileDataWithoutFirstLine := aasmfile.RemoveFirstLine(fileDatString)

	fmt.Println("Removing comments...")

	fileDataWithoutFirstLine = aasmfile.RemoveComments(fileDataWithoutFirstLine)

	fmt.Println("Passing the AASM file to sections parser...")

	firstLevelSections, firstLevelSectionsError := aasmfile.StringToFirstLevelSections(fileDataWithoutFirstLine, 1)

	if firstLevelSectionsError != nil {
		return firstLevelSectionsError
	}

	fmt.Println("Compiling sections tree...")

	sectionsTree, sectionsTreeError := aasmfile.CompileToSectionsTree(firstLevelSections)

	if sectionsTreeError != nil {
		return sectionsTreeError
	}

	return nil
}
