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

type SectionsTree map[string]string

type SectionPolicy struct {
	Required  bool
	Singleton bool
}

var SectionsPolicy = map[string]SectionPolicy{
	"target": {
		Required:  true,
		Singleton: true,
	},
	"namespace": {
		Required:  true,
		Singleton: true,
	},
	"entry": {
		Required:  false,
		Singleton: true,
	},
	"capabilities": {
		Required:  false,
		Singleton: false,
	},
	"batteries": {
		Required:  false,
		Singleton: false,
	},
	"deps": {
		Required:  false,
		Singleton: false,
	},
	"include": {
		Required:  false,
		Singleton: false,
	},
	"exports": {
		Required:  false,
		Singleton: false,
	},
	"constants": {
		Required:  false,
		Singleton: false,
	},
	"reg_aliases": {
		Required:  false,
		Singleton: false,
	},
	"type_aliases": {
		Required:  false,
		Singleton: false,
	},
	"macros": {
		Required:  false,
		Singleton: false,
	},
	"code": {
		Required:  true,
		Singleton: false,
	},
}

func RequiredSections() []string {
	var requiredSections []string

	for sectionName, section := range SectionsPolicy {
		if section.Required {
			requiredSections = append(requiredSections, sectionName)
		}
	}

	return requiredSections
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

func CompileToSectionsTree(sections []Section) (*SectionsTree, error) {
	tree := &SectionsTree{}

	for _, section := range sections {
		policy, policyExists := SectionsPolicy[section.Name]

		if !policyExists {
			return nil, fmt.Errorf("section %s is not a valid section", section.Name)
		}

		record, recordExists := (*tree)[section.Name]

		if recordExists {

			if policy.Singleton {
				return nil, fmt.Errorf("section %s is a singleton section and cannot be duplicated", section.Name)
			}

			(*tree)[section.Name] = fmt.Sprintf("%s\n%s", record, section.Content)

		} else {
			(*tree)[section.Name] = section.Content
		}
	}

	requiredSections := RequiredSections()

	for _, requiredSection := range requiredSections {
		_, sectionExists := (*tree)[requiredSection]

		if !sectionExists {
			return nil, fmt.Errorf("section %s is required but not defined", requiredSection)
		}
	}

	return tree, nil
}
