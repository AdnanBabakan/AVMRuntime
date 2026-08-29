package aasmfile

import "regexp"

func RemoveComments(content string) string {
	re := regexp.MustCompile(`;(.*?)\n`)
	return re.ReplaceAllString(content, "\n")
}
