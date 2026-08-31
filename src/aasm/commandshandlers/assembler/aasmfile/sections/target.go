package sections

import (
	"AVMRuntime/src/aasm/utilities/file"
	"AVMRuntime/src/aasm/utilities/versioning"
	"fmt"
	"slices"
	"strings"
)

type TargetSection struct {
	VM  *versioning.ConstrainedVersion
	ISA *versioning.ConstrainedVersion
}

type TargetPolicy struct {
	Required bool
}

var TargetsPolicies = map[string]TargetPolicy{
	"vm": {
		Required: true,
	},
	"isa": {
		Required: true,
	},
}

func RequiredTargets() []string {
	var requiredTargets []string

	for target, policy := range TargetsPolicies {
		if policy.Required {
			requiredTargets = append(requiredTargets, target)
		}
	}

	return requiredTargets
}

func ToTargetSection(content string) (*TargetSection, error) {
	targetSection := &TargetSection{}

	sanitizedContent := file.RemoveExtraLines(content)

	lineSplit := strings.SplitSeq(sanitizedContent, "\n")

	requiredTargets := RequiredTargets()
	var providedTargets []string

	for line := range lineSplit {
		sanitizedLine := file.SanitizeLine(line)
		fields := strings.Fields(sanitizedLine)
		target := fields[0]
		versionString := fields[1]

		_, policyExists := TargetsPolicies[target]

		if !policyExists {
			return nil, fmt.Errorf("target '%s' is not a valid target", target)
		}

		providedTargets = append(providedTargets, target)

		constrainedVersion, constrainedVersionError := versioning.StringToConstrainedVersion(versionString)

		if constrainedVersionError != nil {
			return nil, constrainedVersionError
		}

		switch target {
		case "vm":
			targetSection.VM = constrainedVersion
		case "isa":
			targetSection.ISA = constrainedVersion
		}
	}

	for _, requiredTarget := range requiredTargets {
		if !slices.Contains(providedTargets, requiredTarget) {
			return nil, fmt.Errorf("target '%s' is required but not provided", requiredTarget)
		}
	}

	return targetSection, nil
}
