package aasmfile

import (
	"AVMRuntime/src/aasm/commandshandlers/assembler/aasmfile/sections"
	"AVMRuntime/src/aasm/utilities/file"
	"fmt"
	"regexp"
	"strings"
)

// raw sections

func RequiredSections() []string {
	var requiredSections []string

	for sectionName, section := range sections.Policy {
		if section.Required {
			requiredSections = append(requiredSections, sectionName)
		}
	}

	return requiredSections
}

func StringToFirstLevelSections(content string, lineOffset int) ([]sections.Section, error) {
	lineSplit := strings.Split(content, "\n")

	sectionNameRe := regexp.MustCompile(`!([0-9a-zA-Z]*)\s`)

	var sectionsSlice = make([]sections.Section, 0)

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

		currentSection := sections.Section{}

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

		sectionsSlice = append(sectionsSlice, currentSection)

		i++
	}

	return sectionsSlice, nil
}

func CompileToSectionsTree(sectionsList []sections.Section) (*sections.Tree, error) {
	tree := &sections.Tree{}

	for _, section := range sectionsList {
		policy, policyExists := sections.Policy[section.Name]

		if !policyExists {
			return nil, fmt.Errorf("section '%s' is not a valid section", section.Name)
		}

		record, recordExists := (*tree)[section.Name]

		if recordExists {

			if policy.Singleton {
				return nil, fmt.Errorf("section '%s' is a singleton section and cannot be duplicated", section.Name)
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
			return nil, fmt.Errorf("section '%s' is required but not defined", requiredSection)
		}
	}

	return tree, nil
}

// full sections

func CompileToFullSectionsTree(sectionsList *sections.Tree) (*sections.FullTree, error) {

	tree := &sections.FullTree{}

	for sectionName, sectionContent := range *sectionsList {
		switch sectionName {
		case "target":
			section, sectionError := sections.ToTargetSection(sectionContent)

			if sectionError != nil {
				return nil, sectionError
			}

			tree.Target = section
		}
	}

	return nil, nil
}
