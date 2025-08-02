package autoenv

import (
	"bufio"
	"os"
	"path/filepath"
	"slices"
)

func (l *Loader) loadEnvFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	f, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			l.options.logger.ErrorF("failed to close file: %s", err)
		}
	}(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		line = trimSpaces(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		line = stripInlineComment(line)
		if len(line) == 0 {
			continue
		}

		if len(line) > 7 && string(line[:7]) == "export " {
			line = trimSpaces(line[7:])
			if len(line) == 0 {
				continue
			}
		}

		i := slices.Index(line, '=')
		if i <= 0 {
			continue
		}
		key := string(trimSpaces(line[:i]))
		val := trimSpaces(line[i+1:])
		if len(val) > 1 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		if err := os.Setenv(key, string(val)); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func trimSpaces(b []byte) []byte {
	start, end := 0, len(b)-1
	for start <= end && (b[start] == ' ' || b[start] == '\t') {
		start++
	}
	for end >= start && (b[end] == ' ' || b[end] == '\t') {
		end--
	}
	if start > end {
		return nil
	}
	return b[start : end+1]
}

func stripInlineComment(line []byte) []byte {
	inSingle := false
	inDouble := false
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return trimSpaces(line[:i])
			}
		}
	}
	return trimSpaces(line)
}
