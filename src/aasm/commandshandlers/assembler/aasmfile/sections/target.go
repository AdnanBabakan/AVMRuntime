package sections

import (
	"AVMRuntime/src/aasm/utilities/versioning"
	"fmt"
)

type TargetSection struct {
	VM  versioning.Version
	ISA versioning.Version
}

func ToTargetSection(content string) (*TargetSection, error) {
	fmt.Println("RUNNING ToTargetSection")

	return nil, nil
}
