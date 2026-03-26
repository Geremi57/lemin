package parser

import (
	"fmt"
	"lem-in/models"
	"lem-in/utils"
)

func ParseFile(filename string) (*models.Colony, error) {
	lines, err := utils.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, cannot read file: %v", err)
	}
	
	if len(lines) == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, empty file")
	}
	
	colony := models.NewColony()
	
	antCount, nextIdx, err := ParseAntCount(lines)
	if err != nil {
		return nil, err
	}
	colony.AntCount = antCount
	
	roomParser := NewRoomParser(colony)
	nextIdx, err = roomParser.ParseRooms(lines, nextIdx)
	if err != nil {
		return nil, err
	}
	
	if colony.Start == nil {
		return nil, fmt.Errorf("ERROR: invalid data format, no start room found")
	}
	if colony.End == nil {
		return nil, fmt.Errorf("ERROR: invalid data format, no end room found")
	}
	
	linkParser := NewLinkParser(colony)
	if err := linkParser.ParseLinks(lines, nextIdx); err != nil {
		return nil, err
	}
	
	return colony, nil
}