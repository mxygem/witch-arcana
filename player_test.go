package witcharcana

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodPlayer(t *testing.T) {
	testCases := []struct {
		name        string
		playerName  string
		clubs       Clubs
		expected    *Player
		expectedErr error
	}{
		{
			name:       "no player found",
			playerName: "Hoeb",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}}},
				},
			},
			expectedErr: fmt.Errorf(`player "Hoeb" not found`),
		},
		{
			name:       "player found",
			playerName: "Quinoa",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
				},
			},
			expected: &Player{Name: "Quinoa", Club: "SP"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.clubs.Player(tc.playerName)

			assert.Equal(t, tc.expected, actual)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlayer(t *testing.T) {
	testCases := []struct {
		name        string
		playerName  string
		clubs       Clubs
		expected    *Player
		expectedPos int
	}{
		{
			name:       "correct index of found player",
			playerName: "Hoeb",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Capsy"}, {Name: "PHTEVEN"}, {Name: "Hoeb"}}},
				},
			},
			expected:    &Player{Name: "Hoeb", Club: "CNT"},
			expectedPos: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pos, actual := player(tc.clubs, tc.playerName)

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedPos, pos)
		})
	}
}

func TestCreatePlayer(t *testing.T) {
	testCases := []struct {
		name        string
		clubName    string
		player      *Player
		clubs       Clubs
		expected    *Player
		expectedErr error
	}{
		{
			name:     "club not found",
			clubName: "DYR",
			player:   &Player{Name: "Andr0meda"},
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
				},
			},
			expectedErr: fmt.Errorf(`club "DYR" not found`),
		},
		{
			name:     "player exists",
			clubName: "EVA",
			player:   &Player{Name: "LadyLuna"},
			clubs: Clubs{
				clubs: map[string]*Club{
					"EVA": {Name: "EVA", Players: []*Player{{Name: "Shotgun"}, {Name: "LadyLuna"}}},
				},
			},
			expectedErr: fmt.Errorf(`creating player: player "LadyLuna" already exists in club "EVA"`),
		},
		{
			name:     "successful create",
			clubName: "BRD",
			player:   &Player{Name: "RedKangaroo"},
			clubs: Clubs{
				clubs: map[string]*Club{
					"BRD": {Name: "BRD", Players: []*Player{{Name: "BlackBad"}, {Name: "MochaGamma"}}},
				},
			},
			expected: &Player{Name: "RedKangaroo", Club: "BRD"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.clubs.CreatePlayer(tc.clubName, tc.player)

			assert.Equal(t, tc.expected, actual)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemovePlayer(t *testing.T) {
	testCases := []struct {
		name        string
		playerName  string
		clubs       Clubs
		expected    Clubs
		expectedErr error
	}{
		{
			name:       "player doesn't exist",
			playerName: "Wishy",
			clubs: Clubs{
				clubs: map[string]*Club{
					"DYR": {Name: "DYR", Players: []*Player{{Name: "RubyBlack"}, {Name: "Spooffy"}}},
				},
			},
			expected: Clubs{
				clubs: map[string]*Club{
					"DYR": {Name: "DYR", Players: []*Player{{Name: "RubyBlack"}, {Name: "Spooffy"}}},
				},
			},
			expectedErr: fmt.Errorf(`removing player: player "Wishy" does not exist`),
		},
		{
			name:       "successfully removed from beginning",
			playerName: "_ScarletRose_",
			clubs: Clubs{
				clubs: map[string]*Club{
					"KMA": {Name: "KMA", Players: []*Player{{Name: "_ScarletRose_"}, {Name: "AriettaRex"}, {Name: "Emeriya"}}},
				},
			},
			expected: Clubs{
				clubs: map[string]*Club{
					"KMA": {Name: "KMA", Players: []*Player{{Name: "AriettaRex"}, {Name: "Emeriya"}}},
				},
			},
		},
		{
			name:       "successfully removed from end",
			playerName: "_ScarletRose_",
			clubs: Clubs{
				clubs: map[string]*Club{
					"KMA": {Name: "KMA", Players: []*Player{{Name: "Emeriya"}, {Name: "Moira"}, {Name: "_ScarletRose_"}}},
				},
			},
			expected: Clubs{
				clubs: map[string]*Club{
					"KMA": {Name: "KMA", Players: []*Player{{Name: "Emeriya"}, {Name: "Moira"}}},
				},
			},
		},
		{
			name:       "successfully removed from end",
			playerName: "_ScarletRose_",
			clubs: Clubs{
				clubs: map[string]*Club{
					"KMA": {Name: "KMA", Players: []*Player{
						{Name: "Emeriya"},
						{Name: "Moira"},
						{Name: "_ScarletRose_"},
						{Name: "AriettaRex"},
						{Name: "ProEvil"},
					}},
				},
			},
			expected: Clubs{
				clubs: map[string]*Club{
					"KMA": {Name: "KMA", Players: []*Player{
						{Name: "Emeriya"},
						{Name: "Moira"},
						{Name: "AriettaRex"},
						{Name: "ProEvil"},
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.clubs.RemovePlayer(tc.playerName)

			assert.Equal(t, tc.expected, tc.clubs)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMovePlayer(t *testing.T) {
	testCases := []struct {
		name          string
		playerName    string
		newClubName   string
		clubs         Clubs
		expected      *Player
		expectedClubs Clubs
		expectedErr   error
	}{
		{
			name:        "player not found",
			playerName:  "treees",
			newClubName: "SP",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
				},
			},
			expectedClubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
				},
			},
			expectedErr: fmt.Errorf(`unable to move player "treees" to "SP": player "treees" does not exist`),
		},
		{
			name:        "new club not found",
			playerName:  "mxygem",
			newClubName: "SP",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				},
			},
			expectedClubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				},
			},
			expectedErr: fmt.Errorf(`unable to move player "mxygem" to "SP": getting new club: no club "SP" found`),
		},
		{
			name:        "player already exists in new club?",
			playerName:  "mxygem",
			newClubName: "SP",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "mxygem"}, {Name: "Jasmin"}}},
				},
			},
			expectedClubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "mxygem"}, {Name: "Jasmin"}}},
				},
			},
			expectedErr: fmt.Errorf(`unable to move player "mxygem" to "SP": creating player in new club: player "mxygem" already exists in club "SP"`),
		},
		{
			name:        "successful move",
			playerName:  "mxygem",
			newClubName: "SP",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
				},
			},
			expected: &Player{Name: "mxygem", Club: "SP"},
			expectedClubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}, {Name: "mxygem"}}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.clubs.MovePlayer(tc.playerName, tc.newClubName)

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedClubs, tc.clubs)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCsBulkUpdatePlayers(t *testing.T) {
	testCases := []struct {
		name        string
		clubs       Clubs
		players     Players
		expected    Clubs
		expectedErr error
	}{
		{
			name: "no player data provided",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
				},
			},
			expected: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
				},
			},
		},
		{
			name: "players in club have levels added",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				},
			},
			players: []*Player{
				{Name: "mxygem", Level: 18, Club: "CNT"},
				{Name: "Hoeb", Level: 15, Club: "CNT"},
			},
			expected: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{
						{Name: "mxygem", Level: 18},
						{Name: "Hoeb", Level: 15},
					}},
				},
			},
		},
		{
			name: "single player added to single club",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "Hoeb"}}},
				},
			},
			players: []*Player{
				{Name: "mxygem", Level: 18, Club: "CNT"},
			},
			expected: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{
						{Name: "Hoeb"},
						{Name: "mxygem", Level: 18},
					}},
				},
			},
		},
		{
			name: "multiple adds and updates across multiple clubs",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "Hoeb"}, {Name: "mxygem", Level: 18}}},
					"MID": {Name: "MID", Players: []*Player{{Name: "AnsaLovesYou"}}},
				},
			},
			players: []*Player{
				{Name: "mxygem", Level: 19, Club: "CNT"},
				{Name: "AnsaLovesYou", Club: "MID", Location: &Location{X: 123, Y: 456}},
				{Name: "Quinoa", Club: "SP", Level: 16},
				{Name: "Jasmine", Club: "SP", Level: 16},
				{Name: "M4rs", Club: "CNT", Level: 15},
			},
			expected: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{
						{Name: "Hoeb"},
						{Name: "mxygem", Level: 19},
						{Name: "M4rs", Level: 15},
					}},
					"MID": {Name: "MID", Players: []*Player{
						{Name: "AnsaLovesYou", Location: &Location{X: 123, Y: 456}},
					}},
					"SP": {Name: "SP", Players: []*Player{
						{Name: "Quinoa", Level: 16},
						{Name: "Jasmine", Level: 16},
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("epl: %d\n", len(tc.expected.clubs["CNT"].Players))
			err := tc.clubs.BulkUpdatePlayers(tc.players)

			assert.Equal(t, tc.expected, tc.clubs)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCsUpdatePlayer(t *testing.T) {
	testCases := []struct {
		name          string
		clubs         Clubs
		player        *Player
		expected      *Player
		expectedClubs Clubs
		expectedErr   error
	}{
		{
			name: "player not found",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{
						{Name: "Hoeb"},
						{Name: "M4rs", Level: 15},
					}},
				},
			},
			player: &Player{
				Name:  "mxygem",
				Club:  "CNT",
				Level: 19,
			},
			expectedErr: fmt.Errorf(`player "mxygem" not found`),
		},
		{
			name: "player found",
			clubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{
						{Name: "Hoeb"},
						{Name: "M4rs", Level: 15},
					}},
				},
			},
			player: &Player{
				Name:  "M4rs",
				Club:  "CNT",
				Level: 16,
			},
			expected: &Player{
				Name:  "M4rs",
				Level: 16,
			},
			expectedClubs: Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{
						{Name: "Hoeb"},
						{Name: "M4rs", Level: 16},
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.clubs.UpdatePlayer(tc.player)

			assert.Equal(t, tc.expected, actual)
			if tc.expectedClubs.clubs != nil {
				assert.Equal(t, tc.expectedClubs, tc.clubs)
			}
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdatePlayer(t *testing.T) {
	testCases := []struct {
		name     string
		op       *Player
		up       *Player
		expected *Player
	}{
		{
			name:     "level",
			op:       &Player{Name: "M4rs", Level: 15},
			up:       &Player{Name: "M4rs", Level: 16, Club: "CNT"},
			expected: &Player{Name: "M4rs", Level: 16},
		},
		{
			name:     "location",
			op:       &Player{Name: "Hoeb", Location: &Location{X: 123, Y: 456}},
			up:       &Player{Name: "Hoeb", Location: &Location{X: 789, Y: 321}},
			expected: &Player{Name: "Hoeb", Location: &Location{X: 789, Y: 321}},
		},
		{
			name:     "in hive",
			op:       &Player{Name: "Quinoa", InHive: false},
			up:       &Player{Name: "Quinoa", InHive: true},
			expected: &Player{Name: "Quinoa", InHive: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, updatePlayer(tc.op, tc.up))
		})
	}
}
