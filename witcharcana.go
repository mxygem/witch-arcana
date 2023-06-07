package witcharcana

// Location stores an X and Y coordinate representing the location of a club or player.
type Location struct {
	X int `json:"x" csv:"x"`
	Y int `json:"y" csv:"y"`
}
