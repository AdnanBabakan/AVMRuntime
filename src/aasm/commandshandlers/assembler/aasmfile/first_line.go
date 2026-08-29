package aasmfile

import (
	"AVMRuntime/src/aasm/utilities/versioning"
	"fmt"
	"strings"
)

type FirstLine struct {
	Tag     string
	Version versioning.Version
}

func (firstLine *FirstLine) VersionString() string {
	return fmt.Sprintf("Major %d - Minor %d - Patch %d - Flag %s", firstLine.Version.Major, firstLine.Version.Minor, firstLine.Version.Patch, firstLine.Version.Flag)
}

func ExtractFirstLine(content string) string {
	firstLine, _, _ := strings.Cut(content, "\n")
	return strings.TrimSuffix(firstLine, "\r")
}

func FirstLineToStruct(line string) (*FirstLine, error) {
	if line == "" {
		return nil, fmt.Errorf("empty first line")
	}

	fields := strings.Fields(line)

	if len(fields) != 2 {
		return nil, fmt.Errorf("first line must contain exactly two fields")
	}

	if fields[0] != "!AASM" {
		return nil, fmt.Errorf("first line must start with !AASM")
	}

	version, err := versioning.StringToVersion(fields[1])

	if err != nil {
		return nil, err
	}

	return &FirstLine{
		Tag:     strings.Trim(fields[0], "!"),
		Version: *version,
	}, nil
}

func ExtractAndCheckFirstLine(line string) (*FirstLine, error) {
	firstLine := ExtractFirstLine(line)
	return FirstLineToStruct(firstLine)
}

func RemoveFirstLine(content string) string {
	_, rest, _ := strings.Cut(content, "\n")
	return strings.TrimSuffix(rest, "\r")
}
