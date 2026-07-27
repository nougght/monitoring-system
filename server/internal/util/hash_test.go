package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CompareHash(t *testing.T) {

	params := defaultHashParams()
	password := "s0m3p@ssw0rd"
	hashString, err := Hash(password, params)
	require.NoError(t, err)

	tests := []struct {
		name               string
		originalHashString string
		password           string
		expectedMatch      bool
		expectedError      error
	}{
		{
			name:               "matching password",
			originalHashString: hashString,
			password:           password,
			expectedMatch:      true,
			expectedError:      nil,
		},
		{
			name:               "incorrect password",
			originalHashString: hashString,
			password:           "3847388437",
			expectedMatch:      false,
			expectedError:      nil,
		},
		{
			name:               "invalid hash string",
			originalHashString: "invalidHashString",
			password:           "3847388437",
			expectedMatch:      false,
			expectedError:      ErrInvalidHashString,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			matching, err := CompareHash(test.originalHashString, test.password)
			assert.Equal(t, test.expectedMatch, matching)
			require.ErrorIs(t, err, test.expectedError)
		})
	}
}
