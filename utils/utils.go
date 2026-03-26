package utils

import (
	"bufio"
	"os"
	"strings"
)

func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lines = append(lines, line)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return lines, nil
}

func IsComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

func IsSpecialCommand(line string) bool {
	return line == "##start" || line == "##end"
}

func IsLink(line string) bool {
	return strings.Contains(line, "-") && !strings.HasPrefix(line, "#")
}