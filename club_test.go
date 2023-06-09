package witcharcana

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClubs(t *testing.T) {
	expected := &Clubs{log: true}

	assert.Equal(t, expected, NewClubs(nil, true))
}

func TestNewClub(t *testing.T) {
	testCases := []struct {
		name     string
		x, y     int
		expected *Club
	}{
		{
			name:     "no location",
			expected: &Club{Name: "no location"},
		},
		{
			name:     "with location",
			x:        123,
			y:        456,
			expected: &Club{Name: "with location", Location: &Location{X: 123, Y: 456}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, NewClub(tc.name, tc.x, tc.y))
		})
	}
}

func TestLoadData(t *testing.T) {
	testCases := []struct {
		name        string
		filename    string
		expected    *Clubs
		expectedErr error
	}{
		{
			name:     "file not found",
			filename: "./testdata/404.json",
			expectedErr: fmt.Errorf("loading club data from file: \"./testdata/404.json\": " +
				"reading: open ./testdata/404.json: no such file or directory"),
		},
		{
			name:     "invalid json",
			filename: "./testdata/invalid.json",
			expectedErr: fmt.Errorf("loading club data from file: \"./testdata/invalid.json\": " +
				"unmarshaling: invalid character '}' looking for beginning of value"),
		},
		{
			name:     "empty file",
			filename: "./testdata/empty.json",
			expected: &Clubs{},
		},
		{
			name:     "successful load",
			filename: "./testdata/clubs.json",
			expected: &Clubs{
				clubs: map[string]*Club{
					"404": {Name: "404", Location: &Location{X: 123, Y: 456}, Players: Players{
						{Name: "DireVoidCat", Level: 15, Might: 51848883},
						{Name: "LoverOnyx", Location: &Location{X: 123, Y: 457}, InHive: true},
					}},

					"AZA": {Name: "AZA", Players: Players{
						{Name: "Fayeee", Level: 18, Might: 70265122, Location: &Location{X: 303, Y: 733}},
						{Name: "Richard"},
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := loadData(tc.filename)

			assert.Equal(t, tc.expected, actual)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCsAll(t *testing.T) {
	testCases := []struct {
		name     string
		clubs    *Clubs
		expected map[string]*Club
	}{
		{
			name: "all clubs returned",
			clubs: &Clubs{
				clubs: map[string]*Club{
					"CNT": {Name: "CNT", Players: []*Player{{Name: "Hoeb"}, {Name: "M4rs"}}},
					"MID": {Name: "MID", Players: []*Player{{Name: "AnsaLovesYou"}, {Name: "Menace"}}},
					"SP":  {Name: "SP", Players: []*Player{{Name: "mxygem"}, {Name: "Quinoa"}, {Name: "Jasmine"}}},
				},
			},
			expected: map[string]*Club{
				"CNT": {Name: "CNT", Players: []*Player{{Name: "Hoeb"}, {Name: "M4rs"}}},
				"MID": {Name: "MID", Players: []*Player{{Name: "AnsaLovesYou"}, {Name: "Menace"}}},
				"SP":  {Name: "SP", Players: []*Player{{Name: "mxygem"}, {Name: "Quinoa"}, {Name: "Jasmine"}}},
			},
		},
		{
			name:     "nil clubs",
			clubs:    &Clubs{},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.clubs.All())
		})
	}
}

func TestCsClub(t *testing.T) {
	testCases := []struct {
		name        string
		clubName    string
		clubs       Clubs
		expected    *Club
		expectedErr error
	}{
		{
			name:     "club not found",
			clubName: "KMA",
			clubs: Clubs{
				clubs: map[string]*Club{
					"SP":  {Name: "SP"},
					"MID": {Name: "MID"},
				}},
			expectedErr: fmt.Errorf(`no club "KMA" found`),
		},
		{
			name:     "club found",
			clubName: "MID",
			clubs: Clubs{
				clubs: map[string]*Club{
					"SP":  {Name: "SP"},
					"MID": {Name: "MID"},
				}},
			expected: &Club{Name: "MID"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.clubs.Club(tc.clubName)

			assert.Equal(t, tc.expected, actual)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCsCreateClub(t *testing.T) {
	testCases := []struct {
		name        string
		club        *Club
		clubs       *Clubs
		expected    *Clubs
		expectedErr error
	}{
		{
			name: "successful creation",
			club: &Club{Name: "CS", Location: &Location{X: 123, Y: 456}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
			expected: &Clubs{
				clubs: map[string]*Club{
					"CS":  {Name: "CS", Location: &Location{X: 123, Y: 456}},
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.clubs.CreateClub(tc.club)

			if tc.expected != nil {
				assert.Equal(t, tc.expected, tc.clubs)
			}
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateClub(t *testing.T) {
	testCases := []struct {
		name        string
		club        *Club
		clubs       *Clubs
		expected    *Clubs
		expectedErr error
	}{
		{
			name: "club already exists",
			club: &Club{Name: "MID"},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
			expectedErr: fmt.Errorf(`club "MID" already exists`),
		},
		{
			name:        "no name given",
			club:        &Club{Location: &Location{X: 123, Y: 456}},
			expectedErr: fmt.Errorf("club name required"),
		},
		{
			name: "successful creation",
			club: &Club{Name: "CS", Location: &Location{X: 123, Y: 456}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
			expected: &Clubs{
				clubs: map[string]*Club{
					"CS":  {Name: "CS", Location: &Location{X: 123, Y: 456}},
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := createClub(tc.clubs, tc.club)

			if tc.expected != nil {
				assert.Equal(t, tc.expected, tc.clubs)
			}
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCsUpdateClub(t *testing.T) {
	testCases := []struct {
		name          string
		updated       *Club
		clubs         *Clubs
		expected      *Club
		expectedClubs *Clubs
		expectedErr   error
	}{
		{
			name:        "no club name in updated info",
			updated:     &Club{Location: &Location{X: 321, Y: 876}},
			expectedErr: fmt.Errorf("updating club: updated club information must contain a name"),
		},
		{
			name:        "invalid x coordinate",
			updated:     &Club{Name: "404", Location: &Location{X: 0, Y: 876}},
			expectedErr: fmt.Errorf("updating club: no location values can be zero. got x: 0 y: 876"),
		},
		{
			name:        "invalid y coordinate",
			updated:     &Club{Name: "404", Location: &Location{X: 123, Y: 0}},
			expectedErr: fmt.Errorf("updating club: no location values can be zero. got x: 123 y: 0"),
		},
		{
			name:    "club doesn't exist",
			updated: &Club{Name: "DYR", Location: &Location{X: 321, Y: 876}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"CS":  {Name: "CS", Location: &Location{X: 123, Y: 456}},
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
			expectedErr: fmt.Errorf(`updating club: no club "DYR" found`),
		},
		{
			name:    "location added",
			updated: &Club{Name: "MID", Location: &Location{X: 321, Y: 876}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
			expected: &Club{Name: "MID", Location: &Location{X: 321, Y: 876}, Players: []*Player{{Name: "Menace"}}},
			expectedClubs: &Clubs{
				clubs: map[string]*Club{
					"MID": {Name: "MID", Location: &Location{X: 321, Y: 876}, Players: []*Player{{Name: "Menace"}}},
				},
			},
		},
		{
			name:    "location updated - only x",
			updated: &Club{Name: "SP", Location: &Location{X: 123, Y: 876}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 321, Y: 876}},
				},
			},
			expected: &Club{Name: "SP", Location: &Location{X: 123, Y: 876}},
			expectedClubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 123, Y: 876}},
				},
			},
		},
		{
			name:    "location updated - only n",
			updated: &Club{Name: "SP", Location: &Location{X: 321, Y: 678}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 321, Y: 876}},
				},
			},
			expected: &Club{Name: "SP", Location: &Location{X: 321, Y: 678}},
			expectedClubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 321, Y: 678}},
				},
			},
		},
		{
			name:    "location updated - both",
			updated: &Club{Name: "SP", Location: &Location{X: 246, Y: 135}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 321, Y: 876}},
				},
			},
			expected: &Club{Name: "SP", Location: &Location{X: 246, Y: 135}},
			expectedClubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 246, Y: 135}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.clubs.UpdateClub(tc.updated)

			assert.Equal(t, tc.expected, actual)
			if tc.expected != nil {
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

func TestUpdateClub(t *testing.T) {
	testCases := []struct {
		name          string
		updated       *Club
		clubs         *Clubs
		expected      *Club
		expectedClubs *Clubs
		expectedErr   error
	}{
		{
			name:        "no club name in updated info",
			updated:     &Club{Location: &Location{X: 321, Y: 876}},
			expectedErr: fmt.Errorf("updated club information must contain a name"),
		},
		{
			name:        "invalid x coordinate",
			updated:     &Club{Name: "404", Location: &Location{X: 0, Y: 876}},
			expectedErr: fmt.Errorf("no location values can be zero. got x: 0 y: 876"),
		},
		{
			name:        "invalid y coordinate",
			updated:     &Club{Name: "404", Location: &Location{X: 123, Y: 0}},
			expectedErr: fmt.Errorf("no location values can be zero. got x: 123 y: 0"),
		},
		{
			name:    "club doesn't exist",
			updated: &Club{Name: "DYR", Location: &Location{X: 321, Y: 876}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"CS":  {Name: "CS", Location: &Location{X: 123, Y: 456}},
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
			expectedErr: fmt.Errorf(`no club "DYR" found`),
		},
		{
			name:    "location added",
			updated: &Club{Name: "MID", Location: &Location{X: 321, Y: 876}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"MID": {Name: "MID", Players: []*Player{{Name: "Menace"}}},
				},
			},
			expected: &Club{Name: "MID", Location: &Location{X: 321, Y: 876}, Players: []*Player{{Name: "Menace"}}},
			expectedClubs: &Clubs{
				clubs: map[string]*Club{
					"MID": {Name: "MID", Location: &Location{X: 321, Y: 876}, Players: []*Player{{Name: "Menace"}}},
				},
			},
		},
		{
			name:    "location updated - only x",
			updated: &Club{Name: "SP", Location: &Location{X: 123, Y: 876}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 321, Y: 876}},
				},
			},
			expected: &Club{Name: "SP", Location: &Location{X: 123, Y: 876}},
			expectedClubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 123, Y: 876}},
				},
			},
		},
		{
			name:    "location updated - only n",
			updated: &Club{Name: "SP", Location: &Location{X: 321, Y: 678}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 321, Y: 876}},
				},
			},
			expected: &Club{Name: "SP", Location: &Location{X: 321, Y: 678}},
			expectedClubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 321, Y: 678}},
				},
			},
		},
		{
			name:    "location updated - both",
			updated: &Club{Name: "SP", Location: &Location{X: 246, Y: 135}},
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 321, Y: 876}},
				},
			},
			expected: &Club{Name: "SP", Location: &Location{X: 246, Y: 135}},
			expectedClubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 246, Y: 135}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := updateClub(tc.clubs, tc.updated)

			assert.Equal(t, tc.expected, actual)
			if tc.expected != nil {
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

func TestCsRemoveClub(t *testing.T) {
	testCases := []struct {
		name        string
		clubName    string
		clubs       *Clubs
		expected    *Clubs
		expectedErr error
	}{
		{
			name:     "club not found",
			clubName: "CCC",
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 246, Y: 135}},
				},
			},
			expectedErr: fmt.Errorf(`removing club: no club "CCC" found`),
		},
		{
			name:     "successful delete",
			clubName: "CCC",
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP":  {Name: "SP", Location: &Location{X: 246, Y: 135}},
					"CCC": {Name: "CCC", Location: &Location{X: 246, Y: 135}, Players: []*Player{{Name: "foo"}}},
				},
			},
			expected: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 246, Y: 135}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cs := *tc.clubs

			err := tc.clubs.RemoveClub(tc.clubName)

			if tc.expectedErr != nil {
				// ensure original data is not changed on err
				assert.Equal(t, &cs, tc.clubs)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.Equal(t, tc.expected, tc.clubs)
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveClub(t *testing.T) {
	testCases := []struct {
		name        string
		clubName    string
		clubs       *Clubs
		expected    *Clubs
		expectedErr error
	}{
		{
			name:     "no name provided",
			clubName: "",
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 246, Y: 135}},
				},
			},
			expectedErr: fmt.Errorf("club name required"),
		},
		{
			name:     "club not found",
			clubName: "CCC",
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 246, Y: 135}},
				},
			},
			expectedErr: fmt.Errorf(`no club "CCC" found`),
		},
		{
			name:     "successful delete",
			clubName: "CCC",
			clubs: &Clubs{
				clubs: map[string]*Club{
					"SP":  {Name: "SP", Location: &Location{X: 246, Y: 135}},
					"CCC": {Name: "CCC", Location: &Location{X: 246, Y: 135}, Players: []*Player{{Name: "foo"}}},
				},
			},
			expected: &Clubs{
				clubs: map[string]*Club{
					"SP": {Name: "SP", Location: &Location{X: 246, Y: 135}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cs := *tc.clubs

			err := removeClub(tc.clubs, tc.clubName)

			if tc.expectedErr != nil {
				// ensure original data is not changed on err
				assert.Equal(t, &cs, tc.clubs)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.Equal(t, tc.expected, tc.clubs)
				assert.NoError(t, err)
			}
		})
	}
}
