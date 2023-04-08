package provider

import "testing"

func TestCompareVersionParts(t *testing.T) {
	type testCase struct {
		name     string
		input    []string
		expected int
	}
	testCases := []testCase{
		{
			name:     "v1 is greater than v2",
			input:    []string{"1", "0"},
			expected: 1,
		},
		{
			name:     "v1 is equal to v2",
			input:    []string{"1", "1"},
			expected: 0,
		},
		{
			name:     "v1 is lesser  than v2",
			input:    []string{"0", "1"},
			expected: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := compareVersionParts(tc.input[0], tc.input[1])
			if actual != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, actual)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	type testCase struct {
		name     string
		input    []string
		expected int
	}
	testCases := []testCase{
		{
			name:     "r1 major is greater than r2 major",
			input:    []string{"v1.0.0", "v0.1.0"},
			expected: 1,
		},
		{
			name:     "r1 minor is greater than r2 minor",
			input:    []string{"v0.2.0", "v0.1.0"},
			expected: 1,
		},
		{
			name:     "r1 patch is greater than r2 patch",
			input:    []string{"v0.1.1", "v0.1.0"},
			expected: 1,
		},
		{
			name:     "r1 equals r2",
			input:    []string{"v0.1.0", "v0.1.0"},
			expected: 0,
		},
		{
			name:     "r1 equals r2",
			input:    []string{"v0.1.0-dev", "v0.1.0-dev"},
			expected: 0,
		},
		{
			name:     "r1 major is lesser than r2 major",
			input:    []string{"v0.1.0", "v1.0.0"},
			expected: -1,
		},
		{
			name:     "r1 minor is lesser than r2 minor",
			input:    []string{"v0.1.0", "v0.2.0"},
			expected: -1,
		},
		{
			name:     "r1 patch is lesser than r2 patch",
			input:    []string{"v0.1.0", "v0.1.1"},
			expected: -1,
		},
		{
			name:     "r1 has pre-release and r2 does not",
			input:    []string{"v0.1.0-alpha", "v0.1.0"},
			expected: -1,
		},
		{
			name:     "r2 has pre-release and r1 does not",
			input:    []string{"v0.1.0", "v0.1.0-alpha"},
			expected: 1,
		},
		{
			name:     "r1 and r2 both have pre-release and are equal",
			input:    []string{"v0.1.0-beta.1", "v0.1.0-beta.1"},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := compareVersions(tc.input[0], tc.input[1])
			if actual != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, actual)
			}
		})
	}
}
