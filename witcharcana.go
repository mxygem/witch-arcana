package witcharcana

// Clubs is a key value store of clubs with keys being the individual club's initials.
type Clubs map[string]*Club

// Club represents a particular club's data and the players that are currently members.
type Club struct {
	Name     string   `json:"name"`
	Location Location `json:"location,omitempty"`
	Players  Players  `json:"players,omitempty"`
}

// Players is a collection of players.
type Players []*Player

// Player represents a player and various data about them.
type Player struct {
	Name     string   `json:"name" csv:"name"`
	Location Location `json:"location,omitempty" csv:"location"`
	InHive   bool     `json:"in_hive,omitempty" csv:"in_hive"`
	Level    int      `json:"level,omitempty" csv:"level,lvl"`
	Might    int64    `json:"might,omitempty" csv:"might"`
	Club     string   `json:"club,omitempty" csv:"club"`
}

// Location stores an X and Y coordinate representing the location of a club or player.
type Location struct {
	X int `json:"x" csv:"x"`
	Y int `json:"y" csv:"y"`
}
