package simulator

import (
	"fmt"
	"lem-in/models"
	"sort"
	"strings"
)

type Ant struct {
	ID       int
	Path     []*models.Room
	Position int
	Finished bool
}

type PathAssignment struct {
	Path     []*models.Room
	AntCount int
}

type Simulator struct {
	colony     *models.Colony
	assignment []PathAssignment
	ants       []*Ant
}

func NewSimulator(colony *models.Colony, assignment []PathAssignment) *Simulator {
	return &Simulator{
		colony:     colony,
		assignment: assignment,
		ants:       make([]*Ant, colony.AntCount),
	}
}

func (s *Simulator) RunSimulation() {
	s.initializeAnts()
	
	finishedCount := 0
	
	for finishedCount < s.colony.AntCount {
		moves := s.processTurn()
		
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
		
		finishedCount = 0
		for _, ant := range s.ants {
			if ant.Finished {
				finishedCount++
			}
		}
	}
}

func (s *Simulator) initializeAnts() {
	antID := 1
	
	for _, assign := range s.assignment {
		for i := 0; i < assign.AntCount; i++ {
			pathCopy := make([]*models.Room, len(assign.Path))
			copy(pathCopy, assign.Path)
			
			s.ants[antID-1] = &Ant{
				ID:       antID,
				Path:     pathCopy,
				Position: 0,
				Finished: false,
			}
			antID++
		}
	}
}

func (s *Simulator) processTurn() []string {
	moves := make([]string, 0)
	
	// Track rooms that will be occupied after this turn (excluding end)
	nextOccupancy := make(map[*models.Room]bool)
	
	// Sort ants by ID for consistent processing
	antsByID := make([]*Ant, len(s.ants))
	copy(antsByID, s.ants)
	sort.Slice(antsByID, func(i, j int) bool {
		return antsByID[i].ID < antsByID[j].ID
	})
	
	// First pass: determine which ants can move
	for _, ant := range antsByID {
		if ant.Finished {
			continue
		}
		
		// Check if at end of path
		if ant.Position >= len(ant.Path)-1 {
			ant.Finished = true
			continue
		}
		
		nextRoom := ant.Path[ant.Position+1]
		
		// Determine if can move
		canMove := false
		
		if nextRoom == s.colony.End {
			// End room always accepts ants
			canMove = true
		} else if nextOccupancy[nextRoom] {
			// Another ant already moving into this room
			canMove = false
		} else if s.isRoomOccupied(nextRoom) {
			// Room currently occupied
			canMove = false
		} else {
			canMove = true
		}
		
		if canMove {
			ant.Position++
			moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom.Name))
			if nextRoom != s.colony.End {
				nextOccupancy[nextRoom] = true
			}
			if nextRoom == s.colony.End {
				ant.Finished = true
			}
		}
	}
	
	// Sort moves by ant ID
	sort.Slice(moves, func(i, j int) bool {
		var id1, id2 int
		fmt.Sscanf(moves[i], "L%d-", &id1)
		fmt.Sscanf(moves[j], "L%d-", &id2)
		return id1 < id2
	})
	
	return moves
}

func (s *Simulator) isRoomOccupied(room *models.Room) bool {
	for _, ant := range s.ants {
		if ant.Finished {
			continue
		}
		if ant.Position < len(ant.Path) && ant.Path[ant.Position] == room {
			return true
		}
	}
	return false
}