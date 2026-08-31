package versioning

import (
	"AVMRuntime/src/aasm/utilities/stringsandchars"
	"fmt"
)

type VersionConstraint int

const (
	LatestVersion VersionConstraint = iota
	AnyVersion
	ExactVersion
	ExceptVersion
	GreaterThanVersion
	GreaterThanOrEqualVersion
	LessThanVersion
	LessThanOrEqualVersion
	MajorCompatibleVersion
	MinorCompatibleVersion
)

type ConstrainedVersion struct {
	Constraint VersionConstraint
	Version    *Version
}

func StringToConstrainedVersion(str string) (*ConstrainedVersion, error) {
	v := &ConstrainedVersion{}

	constraintCharCount := 0

	for _, r := range str {
		if !stringsandchars.IsNum(string(r)) {
			constraintCharCount++
			continue
		}

		break
	}

	constraint := str[0:constraintCharCount]

	switch constraint {
	case "latest", "@":
		v.Constraint = LatestVersion
	case "*":
		v.Constraint = AnyVersion
	case "", "=":
		v.Constraint = ExactVersion
	case "!=", "<>":
		v.Constraint = ExceptVersion
	case ">":
		v.Constraint = GreaterThanVersion
	case ">=":
		v.Constraint = GreaterThanOrEqualVersion
	case "<":
		v.Constraint = LessThanVersion
	case "<=":
		v.Constraint = LessThanOrEqualVersion
	case "^":
		v.Constraint = MajorCompatibleVersion
	case "~":
		v.Constraint = MinorCompatibleVersion
	default:
		return nil, fmt.Errorf("invalid constraint: %s", constraint)
	}

	versionString := str[constraintCharCount:]

	if v.Constraint == LatestVersion || v.Constraint == AnyVersion {
		versionString = "0.0.0"
	}

	version, versionError := StringToVersion(versionString)

	if versionError != nil {
		return nil, versionError
	}

	v.Version = version

	return v, nil
}

func ConstraintToString(constraint VersionConstraint) string {
	switch constraint {
	case LatestVersion:
		return "latest"
	case AnyVersion:
		return "*"
	case ExactVersion:
		return "="
	case ExceptVersion:
		return "!="
	case GreaterThanVersion:
		return ">"
	case GreaterThanOrEqualVersion:
		return ">="
	case LessThanVersion:
		return "<"
	case LessThanOrEqualVersion:
		return "<="
	case MajorCompatibleVersion:
		return "^"
	case MinorCompatibleVersion:
		return "~"
	default:
		panic(fmt.Sprintf("invalid constraint: VersionConstraint(%v)", constraint))
	}
}

func (constrainedVersion *ConstrainedVersion) String() string {
	return fmt.Sprintf("%v%v", ConstraintToString(constrainedVersion.Constraint), constrainedVersion.Version.String())
}

func (constrainedVersion *ConstrainedVersion) DescriptiveString() string {
	switch constrainedVersion.Constraint {
	case LatestVersion:
		return "Latest version"
	case AnyVersion:
		return "Any version"
	case ExactVersion:
		return fmt.Sprintf("Exactly version %s (=%[1]s)", constrainedVersion.Version.String())
	case ExceptVersion:
		return fmt.Sprintf("Except for version %s (!=%[1]s)", constrainedVersion.Version.String())
	case GreaterThanVersion:
		return fmt.Sprintf("Greater than version %s (>%[1]s)", constrainedVersion.Version.String())
	case GreaterThanOrEqualVersion:
		return fmt.Sprintf("Greater than or equal to version %s (>=%[1]s)", constrainedVersion.Version.String())
	case LessThanVersion:
		return fmt.Sprintf("Less than version %s (<%[1]s)", constrainedVersion.Version.String())
	case LessThanOrEqualVersion:
		return fmt.Sprintf("Less than or equal to version %s (<=%[1]s)", constrainedVersion.Version.String())
	case MajorCompatibleVersion:
		fromVersion := *constrainedVersion.Version
		fromVersion.SetMinor(0)
		fromVersion.SetPatch(0)

		toVersion := *constrainedVersion.Version
		toVersion.ApplyToMajor(1)
		toVersion.SetMinor(0)
		toVersion.SetPatch(0)

		return fmt.Sprintf("Major compatible with %s (>=%s && <%s)", constrainedVersion.Version.String(), fromVersion.String(), toVersion.String())
	case MinorCompatibleVersion:
		fromVersion := *constrainedVersion.Version
		fromVersion.SetPatch(0)

		toVersion := *constrainedVersion.Version
		toVersion.ApplyToMinor(1)
		toVersion.SetPatch(0)

		return fmt.Sprintf("Minor compatible with %s (>=%s && <%s)", constrainedVersion.Version.String(), fromVersion.String(), toVersion.String())
	default:
		panic(fmt.Sprintf("invalid constrained version: %v", constrainedVersion))
	}
}
