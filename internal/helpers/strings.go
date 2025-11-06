package helpers

import (
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

func ToSnakeCase(s string) string {
	var result []rune
	var prev rune

	for i, r := range s {
		if r == ' ' || r == '-' {
			if len(result) > 0 && result[len(result)-1] != '_' {
				result = append(result, '_')
			}
			prev = r
			continue
		}

		if unicode.IsUpper(r) {
			lower := unicode.ToLower(r)

			if i > 0 && (unicode.IsLower(prev) || unicode.IsDigit(prev)) {
				if len(result) > 0 && result[len(result)-1] != '_' {
					result = append(result, '_')
				}
			}

			if i > 0 && unicode.IsUpper(prev) {
				if i+1 < len(s) {
					next := rune(s[i+1])
					if unicode.IsLower(next) {
						if len(result) > 0 && result[len(result)-1] != '_' {
							result = append(result, '_')
						}
					}
				}
			}

			result = append(result, lower)
		} else {
			if r == '_' && len(result) > 0 && result[len(result)-1] == '_' {
				continue
			}
			result = append(result, unicode.ToLower(r))
		}

		prev = r
	}

	snake := string(result)
	snake = strings.Trim(snake, "_")
	return snake
}
