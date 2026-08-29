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
