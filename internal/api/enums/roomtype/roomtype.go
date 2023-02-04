package roomtype

// RoomType is the type of a room
type RoomType int

const (
	Private RoomType = iota
	Group
)

func (r RoomType) String() string {
	return [...]string{"private", "group"}[r]
}
