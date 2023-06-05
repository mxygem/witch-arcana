package witcharcana

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClub(t *testing.T) {
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
				"SP":  {Name: "SP"},
				"MID": {Name: "MID"},
			},
			expectedErr: fmt.Errorf(`no club "KMA" found`),
		},
		{
			name:     "club found",
			clubName: "MID",
			clubs: Clubs{
				"SP":  {Name: "SP"},
				"MID": {Name: "MID"},
			},
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
