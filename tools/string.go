package tools

import (
	"strings"
)

func SplitAndTrimSpace(str string, sep string) []string {
	var sList = strings.Split(strings.TrimSpace(str), sep)
	var result = make([]string, 0, len(sList))
	for _, s := range sList {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			result = append(result, s)
		}
	}
	return result
}