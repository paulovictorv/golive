package util

import (
	"strings"
	"unicode"
)

func SplitComma(commaString string) []string {
	split := strings.Split(commaString, ",")

	for i, _ := range split {
		split[i] = strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, split[i])
	}

	return split
}
