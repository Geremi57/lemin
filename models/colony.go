package models

type Room struct {
	Name      string
	X, Y      int
	IsStart   bool
	IsEnd     bool
	Connected []*Room
}

type Colony struct {
	AntCount int
	Rooms    map[string]*Room
	Start    *Room
	End      *Room
}

func NewColony() *Colony {
	return &Colony{
		Rooms: make(map[string]*Room),
	}
}