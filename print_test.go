package witcharcana

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyDiscord(t *testing.T) {
	testCases := []struct {
		name        string
		data        *Clubs
		expected    string
		expectedErr error
	}{
		// 		{
		// 			name: "no data",
		// 			expected: `=====
		// no data
		// =====`,
		// 		},
		{
			name: "no data",
			data: &Clubs{
				clubs: map[string]*Club{
					"404": {Name: "404", Location: &Location{X: 123, Y: 456}, Players: Players{
						{Name: "DireVoidCat", Level: 15, Might: 51848883},
						{Name: "LoverOnyx", Location: &Location{X: 123, Y: 457}, InHive: true},
						{Name: "AnotherPerson", Location: &Location{X: 789, Y: 567}, InHive: true},
					}},

					"AZA": {Name: "AZA", Players: Players{
						{Name: "Fayeee", Level: 18, Might: 70265122, Location: &Location{X: 303, Y: 733}},
						{Name: "Richard"},
					}},
				},
			},
			expected: `==========
Clubs:
	Name: 404
		Loc: 123, 456
		Players:
			Name: DireVoidCat
			Name: LoverOnyx
				Loc: 123, 457
				In Hive: true
			Name: AnotherPerson
				Loc: 789, 567
				In Hive: true
			
	Name: AZA
		Players:
			Name: Fayee
				Loc: 303, 733
				Level: 18
				Might: 70265122
				InHive: False
			Name: Richard

			==========`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer

			err := PrettyDiscord(&b, tc.data)

			fmt.Println(b.String())

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, b.String())
		})
	}
}
