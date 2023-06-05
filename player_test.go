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
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}}},
			},
			expectedErr: fmt.Errorf(`player "Hoeb" not found`),
		},
		{
			name:       "player found",
			playerName: "Quinoa",
			clubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
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
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Capsy"}, {Name: "PHTEVEN"}, {Name: "Hoeb"}}},
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
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
			},
			expectedErr: fmt.Errorf(`club "DYR" not found`),
		},
		{
			name:     "player exists",
			clubName: "EVA",
			player:   &Player{Name: "LadyLuna"},
			clubs: Clubs{
				"EVA": &Club{Name: "EVA", Players: []*Player{{Name: "Shotgun"}, {Name: "LadyLuna"}}},
			},
			expectedErr: fmt.Errorf(`creating player: player "LadyLuna" already exists in club "EVA"`),
		},
		{
			name:     "successful create",
			clubName: "BRD",
			player:   &Player{Name: "RedKangaroo"},
			clubs: Clubs{
				"BRD": &Club{Name: "BRD", Players: []*Player{{Name: "BlackBad"}, {Name: "MochaGamma"}}},
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
				"DYR": &Club{Name: "DYR", Players: []*Player{{Name: "RubyBlack"}, {Name: "Spooffy"}}},
			},
			expected: Clubs{
				"DYR": &Club{Name: "DYR", Players: []*Player{{Name: "RubyBlack"}, {Name: "Spooffy"}}},
			},
			expectedErr: fmt.Errorf(`removing player: player "Wishy" does not exist`),
		},
		{
			name:       "successfully removed from beginning",
			playerName: "_ScarletRose_",
			clubs: Clubs{
				"KMA": &Club{Name: "KMA", Players: []*Player{{Name: "_ScarletRose_"}, {Name: "AriettaRex"}, {Name: "Emeriya"}}},
			},
			expected: Clubs{
				"KMA": &Club{Name: "KMA", Players: []*Player{{Name: "AriettaRex"}, {Name: "Emeriya"}}},
			},
		},
		{
			name:       "successfully removed from end",
			playerName: "_ScarletRose_",
			clubs: Clubs{
				"KMA": &Club{Name: "KMA", Players: []*Player{{Name: "Emeriya"}, {Name: "Moira"}, {Name: "_ScarletRose_"}}},
			},
			expected: Clubs{
				"KMA": &Club{Name: "KMA", Players: []*Player{{Name: "Emeriya"}, {Name: "Moira"}}},
			},
		},
		{
			name:       "successfully removed from end",
			playerName: "_ScarletRose_",
			clubs: Clubs{
				"KMA": &Club{Name: "KMA", Players: []*Player{
					{Name: "Emeriya"},
					{Name: "Moira"},
					{Name: "_ScarletRose_"},
					{Name: "AriettaRex"},
					{Name: "ProEvil"},
				}},
			},
			expected: Clubs{
				"KMA": &Club{Name: "KMA", Players: []*Player{
					{Name: "Emeriya"},
					{Name: "Moira"},
					{Name: "AriettaRex"},
					{Name: "ProEvil"},
				}},
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
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
			},
			expectedClubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
			},
			expectedErr: fmt.Errorf(`unable to move player "treees" to "SP": player "treees" does not exist`),
		},
		{
			name:        "new club not found",
			playerName:  "mxygem",
			newClubName: "SP",
			clubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
			},
			expectedClubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
			},
			expectedErr: fmt.Errorf(`unable to move player "mxygem" to "SP": getting new club: no club "SP" found`),
		},
		{
			name:        "player already exists in new club?",
			playerName:  "mxygem",
			newClubName: "SP",
			clubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "mxygem"}, {Name: "Jasmin"}}},
			},
			expectedClubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "mxygem"}, {Name: "Jasmin"}}},
			},
			expectedErr: fmt.Errorf(`unable to move player "mxygem" to "SP": creating player in new club: player "mxygem" already exists in club "SP"`),
		},
		{
			name:        "successful move",
			playerName:  "mxygem",
			newClubName: "SP",
			clubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
			},
			expected: &Player{Name: "mxygem", Club: "SP"},
			expectedClubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}, {Name: "mxygem"}}},
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

func TestBulkUpdatePlayers(t *testing.T) {
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
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
			},
			expected: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
				"SP":  &Club{Name: "SP", Players: []*Player{{Name: "Quinoa"}, {Name: "Jasmin"}}},
			},
		},
		{
			name: "players in club have levels added",
			clubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "mxygem"}, {Name: "Hoeb"}}},
			},
			players: []*Player{
				{Name: "mxygem", Level: 18, Club: "CNT"},
				{Name: "Hoeb", Level: 15, Club: "CNT"},
			},
			expected: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{
					{Name: "mxygem", Level: 18},
					{Name: "Hoeb", Level: 15},
				}},
			},
		},
		{
			name: "single player added to single club",
			clubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "Hoeb"}}},
			},
			players: []*Player{
				{Name: "mxygem", Level: 18, Club: "CNT"},
			},
			expected: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{
					{Name: "Hoeb"},
					{Name: "mxygem", Level: 18},
				}},
			},
		},
		{
			name: "multiple adds and updates across multiple clubs",
			clubs: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{{Name: "Hoeb"}, {Name: "mxygem", Level: 18}}},
				"MID": &Club{Name: "MID", Players: []*Player{{Name: "AnsaLovesYou"}}},
			},
			players: []*Player{
				{Name: "mxygem", Level: 19, Club: "CNT"},
				{Name: "AnsaLovesYou", Club: "MID", Location: Location{X: 123, Y: 456}},
				{Name: "Quinoa", Club: "SP", Level: 16},
				{Name: "Jasmine", Club: "SP", Level: 16},
				{Name: "M4rs", Club: "CNT", Level: 15},
			},
			expected: Clubs{
				"CNT": &Club{Name: "CNT", Players: []*Player{
					{Name: "Hoeb"},
					{Name: "mxygem", Level: 19},
					{Name: "M4rs", Level: 15},
				}},
				"MID": &Club{Name: "MID", Players: []*Player{
					{Name: "AnsaLovesYou", Location: Location{X: 123, Y: 456}},
				}},
				"SP": &Club{Name: "SP", Players: []*Player{
					{Name: "Quinoa", Level: 16},
					{Name: "Jasmine", Level: 16},
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("epl: %d\n", len(tc.expected["CNT"].Players))
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
