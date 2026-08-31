package stringsandchars

import "strconv"

func IsNum(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
