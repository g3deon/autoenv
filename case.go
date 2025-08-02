package autoenv

const (
	period     = '.'
	underscore = '_'

	asciiUpperA = 'A'
	asciiUpperZ = 'Z'
	asciiLowerA = 'a'
	asciiLowerZ = 'z'
	digit0      = '0'
	digit9      = '9'

	asciiOffset = 32
)

func toSnakeCase(str string) string {
	n := len(str)
	if n == 0 {
		return ""
	}

	out := make([]byte, 0, n+4)
	uppercaseCount := 0

	for i := 0; i < n; i++ {
		curr := str[i]
		next, nextNext := getNextChars(str, i)

		switch {
		case curr == period:
			out = append(out, underscore)
			uppercaseCount = 0

		case isUppercase(curr):
			if shouldAddUnderscoreBeforeUppercase(str, i, uppercaseCount) {
				out = append(out, underscore)
			}

			if isEndOfAcronym(uppercaseCount, next, nextNext) {
				out = append(out, underscore)
				uppercaseCount = 0
			}

			out = append(out, toLower(curr))
			uppercaseCount++

		case isLowercase(curr):
			if isAcronymEnd(uppercaseCount, next) {
				out = append(out, underscore)
			}
			out = append(out, curr)
			uppercaseCount = 0

		default:
			out = append(out, curr)
			uppercaseCount = 0
		}
	}
	return string(out)
}

func getNextChars(str string, i int) (byte, byte) {
	var next, nextNext byte
	if i+1 < len(str) {
		next = str[i+1]
	}
	if i+2 < len(str) {
		nextNext = str[i+2]
	}
	return next, nextNext
}

func isUppercase(c byte) bool {
	return c >= asciiUpperA && c <= asciiUpperZ
}

func isLowercase(c byte) bool {
	return c >= asciiLowerA && c <= asciiLowerZ
}

func isDigit(c byte) bool {
	return c >= digit0 && c <= digit9
}

func toLower(c byte) byte {
	return c + asciiOffset
}

func shouldAddUnderscoreBeforeUppercase(str string, i, uppercaseCount int) bool {
	return uppercaseCount == 0 && i > 0 && (isLowercase(str[i-1]) || isDigit(str[i-1]))
}

func isEndOfAcronym(uppercaseCount int, next, nextNext byte) bool {
	return uppercaseCount >= 2 && isLowercase(next) && !isDigit(nextNext)
}

func isAcronymEnd(uppercaseCount int, next byte) bool {
	return uppercaseCount >= 2 && !isDigit(next)
}
