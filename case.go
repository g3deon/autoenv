package autoenv

import "strings"

func toSnakeCase(str string) string {
	var result strings.Builder

	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}

	return strings.ToUpper(result.String())
}
