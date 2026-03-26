package parser

import (
	"fmt"
	"lem-in/models"
	"lem-in/utils"
	"strings"
)

// LinkParser handles link parsing
type LinkParser struct {
	colony *models.Colony
}

func NewLinkParser(colony *models.Colony) *LinkParser {
	return &LinkParser{
		colony: colony,
	}
}

// ParseLinks parses all links from the remaining lines
func (lp *LinkParser) ParseLinks(lines []string, startIdx int) error {
	for i := startIdx; i < len(lines); i++ {
		line := lines[i]
		
		// Skip empty lines
		if line == "" {
			continue
		}
		
		// Skip comments
		if utils.IsComment(line) && !utils.IsSpecialCommand(line) {
			continue
		}
		
		// Parse link
		if err := lp.parseLink(line); err != nil {
			return err
		}
	}
	
	return nil
}

func (lp *LinkParser) parseLink(line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("ERROR: invalid data format, invalid link: %s", line)
	}
	
	room1Name := strings.TrimSpace(parts[0])
	room2Name := strings.TrimSpace(parts[1])
	
	// Validate rooms exist
	room1, exists1 := lp.colony.Rooms[room1Name]
	room2, exists2 := lp.colony.Rooms[room2Name]
	
	if !exists1 {
		return fmt.Errorf("ERROR: invalid data format, link to unknown room: %s", room1Name)
	}
	if !exists2 {
		return fmt.Errorf("ERROR: invalid data format, link to unknown room: %s", room2Name)
	}
	
	// Prevent self-links
	if room1 == room2 {
		return fmt.Errorf("ERROR: invalid data format, self-link not allowed: %s", line)
	}
	
	// Check for duplicate links
	if lp.isDuplicateLink(room1, room2) {
		return fmt.Errorf("ERROR: invalid data format, duplicate link: %s", line)
	}
	
	// Add bidirectional connection
	room1.Connected = append(room1.Connected, room2)
	room2.Connected = append(room2.Connected, room1)
	
	return nil
}

func (lp *LinkParser) isDuplicateLink(room1, room2 *models.Room) bool {
	for _, conn := range room1.Connected {
		if conn == room2 {
			return true
		}
	}
	return false
}