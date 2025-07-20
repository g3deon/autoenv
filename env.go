package autoenv

import (
	"bufio"
	"os"
	"strings"
)

func loadEnvFile(path string) error {
	exists, err := fileExists(path)
	if err != nil || !exists {
		return err
	}
	return parseFile(path)
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return info != nil, err
}

func parseFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		processLine(scanner.Text())
	}
	return scanner.Err()
}

func processLine(line string) {
	s := strings.TrimSpace(line)
	if s == "" || strings.HasPrefix(s, "#") {
		return
	}
	s = strings.TrimPrefix(s, "export ")
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return
	}
	k := strings.TrimSpace(parts[0])
	v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	os.Setenv(k, v)
}
