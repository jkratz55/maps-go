package maps

import (
	"strconv"
	"strings"
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

func TestSetIfPresent(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		key      string
		val      int
		expected map[string]int
		valueSet bool
	}{
		{
			name: "Value is not Present",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			key: "white",
			val: 10,
			expected: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			valueSet: false,
		},
		{
			name: "Value is Present",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			key: "red",
			val: 10,
			expected: map[string]int{
				"red":    10,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			valueSet: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := SetIfPresent(test.in, test.key, test.val)
			assert.Equal(t, test.valueSet, set)
			assert.Equal(t, test.expected, test.in)
		})
	}
}

func TestSetIfAbsent(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		key      string
		val      int
		expected map[string]int
		valueSet bool
	}{
		{
			name: "Value is not Present",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			key: "white",
			val: 10,
			expected: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
				"white":  10,
			},
			valueSet: true,
		},
		{
			name: "Value is Present",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			key: "red",
			val: 10,
			expected: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			valueSet: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := SetIfAbsent(test.in, test.key, test.val)
			assert.Equal(t, test.valueSet, set)
			assert.Equal(t, test.expected, test.in)
		})
	}
}

func TestClear(t *testing.T) {
	in := map[string]int{
		"red":    1,
		"blue":   2,
		"yellow": 3,
		"green":  4,
		"orange": 5,
	}
	Clear(in)
	assert.Equal(t, 0, len(in))
}

func TestClone(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		expected map[string]int
	}{
		{
			name:     "Nil Map",
			in:       nil,
			expected: nil,
		},
		{
			name: "Non Nil Map",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			expected: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Clone(test.in)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name string
		src  map[string]int
		dst  map[string]int
	}{
		{
			name: "Empty Dest",
			src: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			dst: map[string]int{},
		},
		{
			name: "Existing Dst",
			src: map[string]int{
				"red":    1,
				"blue":   2,
				"yellow": 3,
				"green":  4,
				"orange": 5,
			},
			dst: map[string]int{
				"white":  1,
				"black":  1,
				"green":  1,
				"purple": 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Copy(test.src, test.dst)
			assert.True(t, containsAllElements(test.src, test.dst))
		})
	}
}

func containsAllElements(src, dst map[string]int) bool {
	for key, val := range src {
		if val2, ok := dst[key]; !ok || val != val2 {
			return false
		}
	}
	return true
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		left     map[string]int
		right    map[string]int
		expected bool
	}{
		{
			name: "Not Equal - Different Len",
			left: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
			},
			right: map[string]int{
				"red":    1,
				"blue":   2,
				"white":  3,
				"orange": 4,
			},
			expected: false,
		},
		{
			name: "Not Equal - Same Len",
			left: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
			},
			right: map[string]int{
				"red":   1,
				"blue":  2,
				"white": 3,
			},
			expected: false,
		},
		{
			name: "Equal",
			left: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
			},
			right: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, Equal(test.left, test.right))
		})
	}
}

func TestEntries(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		expected []Entry[string, int]
	}{
		{
			name:     "Nil Map",
			in:       nil,
			expected: []Entry[string, int]{},
		},
		{
			name: "Populated Map",
			in: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
			},
			expected: []Entry[string, int]{
				{
					Key:   "red",
					Value: 1,
				},
				{
					Key:   "blue",
					Value: 2,
				},
				{
					Key:   "green",
					Value: 3,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.ElementsMatch(t, test.expected, Entries(test.in))
		})
	}
}

func TestMapEntries(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		mapper   EntryMapper[string, string, int, string]
		expected map[string]string
	}{
		{
			name: "Transform Key and Value",
			in: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
			},
			mapper: func(key string, val int) (string, string) {
				newVal := strconv.Itoa(val)
				return strings.ToUpper(key), newVal
			},
			expected: map[string]string{
				"RED":   "1",
				"BLUE":  "2",
				"GREEN": "3",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := MapEntries[map[string]int, map[string]string](test.in, test.mapper)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		in        map[string]int
		predicate Predicate[string, int]
		expected  map[string]int
	}{
		{
			name: "Filter Keys That Don't Contain 'R'",
			in: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
			},
			predicate: func(key string, val int) bool {
				return strings.Contains(key, "r")
			},
			expected: map[string]int{
				"red":   1,
				"green": 3,
			},
		},
		{
			name: "Filter Out Odd Number Values",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"green":  3,
				"orange": 4,
			},
			predicate: func(key string, val int) bool {
				return val%2 == 0
			},
			expected: map[string]int{
				"blue":   2,
				"orange": 4,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, Filter(test.in, test.predicate))
		})
	}
}

func TestTakeIf(t *testing.T) {

	var invoked int

	tests := []struct {
		name         string
		in           map[string]int
		pred         Predicate[string, int]
		fn           func(k string, v int)
		invokedCount int
	}{
		{
			name: "No Predicate Matches",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"green":  3,
				"orange": 4,
			},
			pred: func(key string, val int) bool {
				if val > 10 {
					return true
				}
				return false
			},
			fn: func(k string, v int) {
				invoked++
			},
			invokedCount: 0,
		},
		{
			name: "Matches 2 Entries",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"green":  3,
				"orange": 4,
			},
			pred: func(key string, val int) bool {
				return val%2 == 0
			},
			fn: func(k string, v int) {
				invoked++
			},
			invokedCount: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			invoked = 0
			TakeIf(test.in, test.pred, test.fn)
			assert.Equal(t, test.invokedCount, invoked)
		})
	}
}

func TestDiff(t *testing.T) {
	m1 := map[string]int{
		"red":   1,
		"blue":  2,
		"green": 3,
		"white": 4,
	}
	m2 := map[string]int{
		"red":   1,
		"blue":  1,
		"green": 3,
		"black": 4,
	}

	actual := Diff(m1, m2)

	assert.Equal(t, 3, len(actual))

	val, ok := actual["blue"]
	assert.True(t, ok)
	assert.Equal(t, 2, val.Left)
	assert.Equal(t, 1, val.Right)
	assert.Equal(t, DiffValue, val.Reason)

	val, ok = actual["white"]
	assert.True(t, ok)
	assert.Equal(t, 4, val.Left)
	assert.Equal(t, 0, val.Right)
	assert.Equal(t, DiffMissingRight, val.Reason)

	val, ok = actual["black"]
	assert.True(t, ok)
	assert.Equal(t, 0, val.Left)
	assert.Equal(t, 4, val.Right)
	assert.Equal(t, DiffMissingLeft, val.Reason)
}

func TestMapToSlice(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		mapper   func(k string, v int) string
		expected []string
	}{
		{
			name: "Map to []string",
			in: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
			},
			mapper: func(k string, v int) string {
				return k + "|" + strconv.Itoa(v)
			},
			expected: []string{
				"red|1",
				"blue|2",
				"green|3",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := MapToSlice(test.in, test.mapper)
			assert.ElementsMatch(t, test.expected, actual)
		})
	}
}

func TestKeyDiff(t *testing.T) {
	tests := []struct {
		name          string
		left          map[string]int
		right         map[string]int
		expectedLeft  []string
		expectedRight []string
	}{
		{
			name: "Both Have Additional Keys & Values",
			left: map[string]int{
				"red":    1,
				"blue":   2,
				"green":  3,
				"orange": 4,
			},
			right: map[string]int{
				"red":   1,
				"blue":  2,
				"green": 3,
				"pink":  4,
			},
			expectedLeft:  []string{"orange"},
			expectedRight: []string{"pink"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualLeft, actualRight := KeyDiff(test.left, test.right)
			assert.ElementsMatch(t, test.expectedLeft, actualLeft)
			assert.ElementsMatch(t, test.expectedRight, actualRight)
		})
	}
}

func TestInvert(t *testing.T) {
	tests := []struct {
		name     string
		in       map[string]int
		expected map[int]string
	}{
		{
			name: "No Collisions on Invert",
			in: map[string]int{
				"red":    1,
				"blue":   2,
				"green":  3,
				"orange": 4,
			},
			expected: map[int]string{
				1: "red",
				2: "blue",
				3: "green",
				4: "orange",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Invert(test.in)
			assert.Equal(t, test.expected, actual)
		})
	}
}
