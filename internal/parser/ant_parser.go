package parser

import (
	"fmt"
	"lem-in/utils"
	"strconv"
	"strings"
)

func ParseAntCount(lines []string) (int, int, error) {
	for idx, line := range lines {
		line = strings.TrimSpace(line)

		//skip empty line
		if line == "" {
			continue
		}

		//skip normal comments
		if strings.HasPrefix(line, "#") && !utils.IsSpecialCommand(line) {
			continue
		}

		//try to parse
		ants, err := strconv.Atoi(line)
		if err != nil {
			return 0, 0, fmt.Errorf("ERROR: invalid data format, invalid number of Ants")
		}
		// return ants + next index to start parsing rooms
		return ants, idx + 1, nil
	}
	return 0, 0, fmt.Errorf("ERROR: could not parse ants, no count found")
}
