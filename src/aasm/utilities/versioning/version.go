package versioning

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Flag  string
}

func StringToVersion(str string) (*Version, error) {
	versionAndFlag := strings.Split(str, "-")

	result := &Version{}

	if len(versionAndFlag) == 2 {
		result.Flag = versionAndFlag[1]
	} else {
		result.Flag = "stable"
	}

	versionParts := strings.Split(versionAndFlag[0], ".")

	if len(versionParts) != 3 {
		return nil, fmt.Errorf("invalid version format: %s", versionAndFlag[0])
	}

	result.Major, _ = strconv.Atoi(versionParts[0])
	result.Minor, _ = strconv.Atoi(versionParts[1])
	result.Patch, _ = strconv.Atoi(versionParts[2])

	return result, nil
}

func (version *Version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", version.Major, version.Minor, version.Patch, version.Flag)
}

func (version *Version) SetMajor(major int) {
	version.Major = major
}

func (version *Version) ApplyToMajor(offset int) {
	version.Major += offset
}

func (version *Version) SetMinor(minor int) {
	version.Minor = minor
}

func (version *Version) ApplyToMinor(offset int) {
	version.Minor += offset
}

func (version *Version) SetPatch(patch int) {
	version.Patch = patch
}

func (version *Version) ApplyToPatch(offset int) {
	version.Patch += offset
}

func (version *Version) SetFlags(flag string) {
	version.Flag = flag
}

func (version *Version) DoesComplyWith(constrainedVersion *ConstrainedVersion) bool {
	switch constrainedVersion.Constraint {
	case LatestVersion:
		// TODO: later, I have check whether the latest version is fetched, but for now any version matches the latest
		return true
	case AnyVersion:
		return false
	case ExactVersion:
		return version.Major == constrainedVersion.Version.Major &&
			version.Minor == constrainedVersion.Version.Minor &&
			version.Patch == constrainedVersion.Version.Patch
	case ExceptVersion:
		return version.Major != constrainedVersion.Version.Major ||
			version.Minor != constrainedVersion.Version.Minor ||
			version.Patch != constrainedVersion.Version.Patch
	case GreaterThanVersion:
		if version.Major > constrainedVersion.Version.Major {
			return true
		} else if version.Major == constrainedVersion.Version.Major {
			if version.Minor > constrainedVersion.Version.Minor {
				return true
			} else if version.Minor == constrainedVersion.Version.Minor {
				if version.Patch > constrainedVersion.Version.Patch {
					return true
				}
				return false
			}
			return false
		}
		return false
	case GreaterThanOrEqualVersion:
		if version.Major > constrainedVersion.Version.Major {
			return true
		} else if version.Major == constrainedVersion.Version.Major {
			if version.Minor > constrainedVersion.Version.Minor {
				return true
			} else if version.Minor == constrainedVersion.Version.Minor {
				if version.Patch >= constrainedVersion.Version.Patch {
					return true
				}
				return false
			}
			return false
		}
		return false
	case LessThanVersion:
		if version.Major < constrainedVersion.Version.Major {
			return true
		} else if version.Major == constrainedVersion.Version.Major {
			if version.Minor < constrainedVersion.Version.Minor {
				return true
			} else if version.Minor == constrainedVersion.Version.Minor {
				if version.Patch < constrainedVersion.Version.Patch {
					return true
				}
				return false
			}
			return false
		}
		return false
	case LessThanOrEqualVersion:
		if version.Major < constrainedVersion.Version.Major {
			return true
		} else if version.Major == constrainedVersion.Version.Major {
			if version.Minor < constrainedVersion.Version.Minor {
				return true
			} else if version.Minor == constrainedVersion.Version.Minor {
				if version.Patch <= constrainedVersion.Version.Patch {
					return true
				}
				return false
			}
			return false
		}
		return false
	case MajorCompatibleVersion:
		fromVersion := *constrainedVersion.Version
		fromVersion.SetMinor(0)
		fromVersion.SetPatch(0)

		toVersion := *constrainedVersion.Version
		toVersion.ApplyToMajor(1)
		toVersion.SetMinor(0)
		toVersion.SetPatch(0)

		return version.Major >= fromVersion.Major && version.Major < toVersion.Major
	case MinorCompatibleVersion:
		if version.Major != constrainedVersion.Version.Major {
			return false
		}

		fromVersion := *constrainedVersion.Version
		fromVersion.SetPatch(0)

		toVersion := *constrainedVersion.Version
		toVersion.ApplyToMinor(1)
		toVersion.SetPatch(0)

		return version.Minor >= fromVersion.Minor && version.Minor < toVersion.Minor
	default:
		return false
	}
}
