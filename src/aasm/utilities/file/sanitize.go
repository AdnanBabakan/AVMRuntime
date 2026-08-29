package file

import (
	"regexp"
	"strings"
)

func RemoveExtraLines(content string) string {
	s := strings.TrimSpace(content)
	s = strings.Replace(s, "\r\n", "\n", -1)
	s = strings.Replace(s, "\r", "\n", -1)

	re := regexp.MustCompile(`\n{2,}`)
	s = re.ReplaceAllString(s, "\n")

	return s
}

func SanitizeLine(content string) string {
	s := strings.TrimSpace(content)
	s = strings.Replace(s, "\r\n", "\n", -1)
	s = strings.Replace(s, "\r", "\n", -1)

	re := regexp.MustCompile(`\n{2,}`)
	s = re.ReplaceAllString(s, "\n")

	return s
}
