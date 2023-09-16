package interval_test

import (
	interval "github.com/marshauf/coding_challenges/interval_merge/go"
	"testing"
)

type row struct {
	name     string
	input    []interval.Interval
	expected []interval.Interval
}

func equal(a, b []interval.Interval) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestMerge(t *testing.T) {
	table := []row{
		{
			name:     "Empty",
			input:    []interval.Interval{},
			expected: []interval.Interval{},
		}, {
			name:     "One",
			input:    []interval.Interval{interval.New(0, 1)},
			expected: []interval.Interval{interval.New(0, 1)},
		}, {
			name:     "intervals which do not overlap",
			input:    []interval.Interval{interval.New(0, 1), interval.New(2, 3), interval.New(5, 12)},
			expected: []interval.Interval{interval.New(0, 1), interval.New(2, 3), interval.New(5, 12)},
		}, {
			name:     "example",
			input:    []interval.Interval{interval.New(25, 30), interval.New(2, 19), interval.New(14, 23), interval.New(4, 8)},
			expected: []interval.Interval{interval.New(2, 23), interval.New(25, 30)},
		}, {
			name:     "example with an interval at the end, which includes all intervals before it",
			input:    []interval.Interval{interval.New(25, 30), interval.New(2, 19), interval.New(14, 23), interval.New(4, 8), interval.New(0, 40)},
			expected: []interval.Interval{interval.New(0, 40)},
		}, {
			name:     "negative ranges",
			input:    []interval.Interval{interval.New(-4, 0), interval.New(-1, 2)},
			expected: []interval.Interval{interval.New(-4, 2)},
		},
	}

	for _, row := range table {
		t.Run(row.name, func(t *testing.T) {
			res := interval.Merge(row.input)

			if !equal(res, row.expected) {
				t.Errorf("Expected %v to be merged to %v, got %v", row.input, row.expected, res)
			}

		})
	}
}
