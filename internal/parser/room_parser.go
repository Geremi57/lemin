package parser

import (
	"fmt"
	"lem-in/models"
	"lem-in/utils"
	"strconv"
	"strings"
)

type RoomParser struct {
	colony *models.Colony
}

func NewRoomParser(c *models.Colony) *RoomParser {
	return &RoomParser{colony: c}
}

func (rp *RoomParser) ParseRooms(lines []string, startIdx int) (int, error) {
	isStart := false
	isEnd := false

	for idx := startIdx; idx < len(lines); idx++ {
		line := strings.TrimSpace(lines[idx])

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") && !utils.IsSpecialCommand(line) {
			continue
		}

		//detect special command
		if utils.IsSpecialCommand(line) {
			if line == "##start" {
				isStart = true
			} else {
				isEnd = true
			}
			continue
		}

		// check if this line looks like a link instead of a room
		if utils.IsLink(line) {
			return idx, nil
		}

		//parse room
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			return 0, fmt.Errorf("ERROR: invalid data format")
		}

		name := parts[0]
		x, err1 := strconv.Atoi(parts[1])
		y, err2 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil {
			return 0, fmt.Errorf("ERROR: invalid data format")
		}

		//validate name
		if strings.HasPrefix(name, "L") || utils.IsComment(name) {
			return 0, fmt.Errorf("ERROR: invalid data format")

		}

		//check duplicates, use comma , ok method
		if _, exists := rp.colony.Rooms[name]; exists {
			return 0, fmt.Errorf("ERROR: duplicate room names")

		}

		//create room
		room := &models.Room{Name: name, X: x, Y: y}

		rp.colony.Rooms[name] = room

		//assign start and end
		if isStart {
			rp.colony.Start = room
			isStart = false
		} else if isEnd {
			rp.colony.End = room
			isEnd = false

		}
	}
	return len(lines), nil

}
