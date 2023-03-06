package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {

	tests := []struct {
		name     string
		left     map[string]int
		right    map[string]int
		resolver ConflictResolver[int]
		expected map[string]int
	}{
		{
			name: "No Conflicts",
			left: map[string]int{
				"white": 3,
				"black": 1,
				"green": 5,
			},
			right: map[string]int{
				"red":    7,
				"blue":   4,
				"orange": 11,
			},
			resolver: NopResolver[int](),
			expected: map[string]int{
				"white":  3,
				"black":  1,
				"green":  5,
				"red":    7,
				"blue":   4,
				"orange": 11,
			},
		},
		{
			name: "Has Conflicts",
			left: map[string]int{
				"white": 3,
				"black": 1,
				"green": 5,
			},
			right: map[string]int{
				"red":    7,
				"blue":   4,
				"orange": 11,
				"white":  5,
				"green":  5,
			},
			resolver: func(left, right int) int {
				return left + right
			},
			expected: map[string]int{
				"white":  8,
				"black":  1,
				"green":  10,
				"red":    7,
				"blue":   4,
				"orange": 11,
			},
		},
	}

	for _, test := range tests {
		actual := Merge(test.resolver, test.left, test.right)
		assert.Equal(t, test.expected, actual)
	}
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		expected []string
	}{
		{
			name:     "Empty Map",
			in:       map[string]int{},
			expected: []string{},
		},
		{
			name: "Has Entries",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			expected: []string{"red", "blue", "yellow", "green", "orange"},
		},
	}

	for _, test := range tests {
		actual := Keys(test.in)
		assert.ElementsMatch(t, test.expected, actual)
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		expected []int
	}{
		{
			name:     "Empty Map",
			in:       map[string]int{},
			expected: []int{},
		},
		{
			name: "Has Entries",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Values(test.in)
			assert.ElementsMatch(t, test.expected, actual)
		})
	}
}

func TestGetOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		key      string
		fallback int
		expected int
	}{
		{
			name: "Contains Key",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			key:      "red",
			fallback: 10,
			expected: 1,
		},
		{
			name: "Does Not Contain Key",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			key:      "gray",
			fallback: 10,
			expected: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := GetOrDefault(test.in, test.key, test.fallback)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestGetOrPanic(t *testing.T) {
	tests := []struct {
		name        string
		in          map[string]int
		key         string
		shouldPanic bool
		expected    int
	}{
		{
			name: "Key Doesn't Exist",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			key:         "gray",
			shouldPanic: true,
			expected:    0,
		},
		{
			name: "Key Exist",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			key:         "yellow",
			shouldPanic: false,
			expected:    3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() {
					_ = GetOrPanic(test.in, test.key)
				})
			} else {
				assert.NotPanics(t, func() {
					actual := GetOrPanic(test.in, test.key)
					assert.Equal(t, test.expected, actual)
				})
			}
		})
	}
}
