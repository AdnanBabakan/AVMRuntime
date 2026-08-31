package aasmfile

import (
	"AVMRuntime/src/aasm/utilities/file"
	"fmt"
	"regexp"
	"strings"
)

type Section struct {
	Name    string
	Content string
}

func StringToFirstLevelSections(content string, lineOffset int) ([]Section, error) {
	lineSplit := strings.Split(content, "\n")

	sectionNameRe := regexp.MustCompile(`!([0-9a-zA-Z]*)\s`)

	var sections = make([]Section, 0)

	for i := 0; i < len(lineSplit); {
		line := file.SanitizeLine(lineSplit[i])
		physicalLineIndex := lineOffset + i

		if line == "" {
			i++
			continue
		}

		if !strings.HasPrefix(line, "!") {
			return nil, fmt.Errorf("illegal section definition at line %d: %s", physicalLineIndex, line)
		}

		currentSection := Section{}

		sectionName := strings.Replace(strings.TrimSpace(sectionNameRe.FindString(line)), "!", "", 1)

		if sectionName == "" {
			return nil, fmt.Errorf("illegal section definition at line %d: %s", physicalLineIndex, line)
		}

		currentSection.Name = sectionName

		sectionContent := ""

		if !strings.Contains(line, "{") {
			// single-line section

			sectionContent = strings.TrimSpace(sectionNameRe.ReplaceAllString(line, ""))

		} else {
			// multi-line section

			blockOpeningsCount := 1

			j := i + 1

			for j < len(lineSplit) {
				innerLine := file.SanitizeLine(lineSplit[j])

				if blockOpeningsCount == 0 {
					break
				}

				if innerLine == "}" {
					blockOpeningsCount--
					continue
				}

				if strings.Contains(innerLine, "{") {
					blockOpeningsCount++
				}

				sectionContent = fmt.Sprintf("%s%s\n", sectionContent, innerLine)

				j++
			}

			i += j
		}

		currentSection.Content = sectionContent

		sections = append(sections, currentSection)

		i++
	}

	return sections, nil
}
