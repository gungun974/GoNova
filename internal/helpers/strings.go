package helpers

import (
	"regexp"
	"strings"
	"unicode"
)

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func LowerFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

func ToSnakeCase(input string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := re.ReplaceAllString(input, `${1}_${2}`)

	return strings.ToLower(snake)
}
