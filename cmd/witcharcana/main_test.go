package main

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestBulkUpdatePlayers(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		cs       clubs
// 		ps       players
// 		expected clubs
// 	}{
// 		{
// 			name:     "no players in",
// 			cs:       clubs{"c0": {Name: "c0", Location: location{X: 123, Y: 456}, Players: players{}}},
// 			expected: clubs{"c0": {Name: "c0", Location: location{X: 123, Y: 456}, Players: players{}}},
// 		},
// 		{
// 			name: "add new club and new player",
// 			cs:   clubs{},
// 			ps: players{
// 				testPlayer(),
// 			},
// 			expected: clubs{
// 				"c0": {
// 					Name: "c0",
// 					Players: players{
// 						testPlayer(),
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "add new club and new player, plus additional player",
// 			cs:   clubs{},
// 			ps: players{
// 				&player{Name: "p0", Club: "c0"},
// 				&player{Name: "p1", Club: "c0"},
// 			},
// 			expected: clubs{
// 				"c0": {
// 					Name: "c0",
// 					Players: players{
// 						&player{Name: "p0", Club: "c0"},
// 						&player{Name: "p1", Club: "c0"},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "club with players exists, add new",
// 			cs: clubs{
// 				"c0": {
// 					Name: "c0",
// 					Players: players{
// 						&player{Name: "p0", Club: "c0"},
// 					},
// 				},
// 			},
// 			ps: players{
// 				&player{Name: "p100", Club: "c0"},
// 			},
// 			expected: clubs{
// 				"c0": {
// 					Name: "c0",
// 					Players: players{
// 						&player{Name: "p0", Club: "c0"},
// 						&player{Name: "p100", Club: "c0"},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "update existing player",
// 			cs: clubs{
// 				"c0": {
// 					Name: "c0",
// 					Players: players{
// 						&player{Name: "p0", Club: "c0", Level: 10},
// 					},
// 				},
// 			},
// 			ps: players{
// 				&player{Name: "p0", Club: "c0", Level: 20},
// 			},
// 			expected: clubs{
// 				"c0": {
// 					Name: "c0",
// 					Players: players{
// 						&player{Name: "p0", Club: "c0", Level: 20},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "add new players and then update when dupe found",
// 			cs: clubs{
// 				"c0": {
// 					Name:    "c0",
// 					Players: players{},
// 				},
// 			},
// 			ps: players{
// 				&player{Name: "p0", Club: "c0", Level: 10},
// 				&player{Name: "p1", Club: "c0", Level: 11},
// 				&player{Name: "p2", Club: "c0", Level: 12},
// 				&player{Name: "p0", Club: "c0", Level: 13},
// 			},
// 			expected: clubs{
// 				"c0": {
// 					Name: "c0",
// 					Players: players{
// 						&player{Name: "p0", Club: "c0", Level: 13},
// 						&player{Name: "p1", Club: "c0", Level: 11},
// 						&player{Name: "p2", Club: "c0", Level: 12},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "combination of adding, updating players, and new clubs",
// 			cs: clubs{
// 				"c0-no-players": {
// 					Name:    "c0-no-players",
// 					Players: players{},
// 				},
// 				"c1-existing": {
// 					Name: "c1-existing",
// 					Players: players{
// 						&player{Name: "p0-c1", Club: "c1-existing", Level: 10},
// 					},
// 				},
// 			},
// 			ps: players{
// 				&player{Name: "p0-c1", Club: "c1-existing", Level: 18, Location: location{X: 111, Y: 111}},
// 				&player{Name: "p2-c2", Club: "c2-new-club", Location: location{X: 222, Y: 222}},
// 				&player{Name: "p10-c0", Club: "c0-no-players", Level: 12},
// 				&player{Name: "p1-c2", Club: "c2-new-club", Level: 13, InHive: true},
// 				&player{Name: "p2-c2", Club: "c2-new-club", Level: 20, Location: location{X: 333, Y: 333}, InHive: true},
// 				&player{Name: "p1-c1", Club: "c1-existing"},
// 			},
// 			expected: clubs{
// 				"c0-no-players": {
// 					Name: "c0-no-players",
// 					Players: players{
// 						&player{Name: "p10-c0", Club: "c0-no-players", Level: 12},
// 					},
// 				},
// 				"c1-existing": {
// 					Name: "c1-existing",
// 					Players: players{
// 						&player{Name: "p0-c1", Club: "c1-existing", Level: 18, Location: location{X: 111, Y: 111}},
// 						&player{Name: "p1-c1", Club: "c1-existing"},
// 					},
// 				},
// 				"c2-new-club": {
// 					Name: "c2-new-club",
// 					Players: players{
// 						&player{Name: "p2-c2", Club: "c2-new-club", Level: 20, Location: location{X: 333, Y: 333}, InHive: true},
// 						&player{Name: "p1-c2", Club: "c2-new-club", Level: 13, InHive: true},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			bulkUpdatePlayers(tc.cs, tc.ps)

// 			assert.Equal(t, tc.expected, tc.cs)
// 		})
// 	}
// }

// func TestUpdatePlayerClubKnown(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		club     *club
// 		player   *player
// 		expected *club
// 	}{
// 		{
// 			name: "nil player",
// 			club: &club{
// 				Name: "c400",
// 				Players: players{
// 					testPlayerSpec("c0", "p0"),
// 				},
// 			},
// 			expected: &club{
// 				Name: "c400",
// 				Players: players{
// 					testPlayerSpec("c0", "p0"),
// 				},
// 			},
// 		},
// 		{
// 			name: "nil player",
// 			club: &club{
// 				Name: "c404",
// 				Players: players{
// 					testPlayerSpec("c0", "p0"),
// 				},
// 			},
// 			player: testPlayerSpec("c0", "p404"),
// 			expected: &club{
// 				Name: "c404",
// 				Players: players{
// 					testPlayerSpec("c0", "p0"),
// 				},
// 			},
// 		},
// 		{
// 			name: "match - single player",
// 			club: &club{
// 				Name: "c0",
// 				Players: players{
// 					&player{
// 						Name:     "p0",
// 						Club:     "c0",
// 						Level:    15,
// 						Location: location{X: 321, Y: 987},
// 					},
// 				},
// 			},
// 			player: &player{
// 				Name:     "p0",
// 				Level:    16,
// 				Location: location{X: 333, Y: 666},
// 			},
// 			expected: &club{
// 				Name: "c0",
// 				Players: players{
// 					&player{
// 						Name:     "p0",
// 						Club:     "c0",
// 						Level:    16,
// 						Location: location{X: 333, Y: 666},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "match - end of players",
// 			club: &club{
// 				Name: "c0",
// 				Players: players{
// 					testPlayerSpec("c0", "p0"),
// 					testPlayerSpec("c0", "p1"),
// 					testPlayerSpec("c0", "p2"),
// 					testPlayerSpec("c0", "p3"),
// 					testPlayerSpec("c0", "p4"),
// 					testPlayerSpec("c0", "p5"),
// 				},
// 			},
// 			player: &player{
// 				Name:     "p5",
// 				Level:    20,
// 				Location: location{X: 111, Y: 222},
// 			},
// 			expected: &club{
// 				Name: "c0",
// 				Players: players{
// 					testPlayerSpec("c0", "p0"),
// 					testPlayerSpec("c0", "p1"),
// 					testPlayerSpec("c0", "p2"),
// 					testPlayerSpec("c0", "p3"),
// 					testPlayerSpec("c0", "p4"),
// 					&player{Name: "p5", Club: "c0", Level: 20, Location: location{X: 111, Y: 222}},
// 				},
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			updatePlayerClubKnown(tc.club, tc.player)

// 			assert.Equal(t, tc.expected, tc.club)
// 		})
// 	}
// }

// func testPlayer() *player {
// 	return &player{
// 		Name:     "p0",
// 		Location: location{X: 123, Y: 789},
// 		Level:    15,
// 		Club:     "c0",
// 	}
// }

// func testPlayerSpec(cid, pid string) *player {
// 	return &player{
// 		Name:     pid,
// 		Location: location{X: 123, Y: 789},
// 		Level:    15,
// 		Club:     cid,
// 	}
// }
