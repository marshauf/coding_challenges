package interval_test

import (
	"fmt"
	"math/rand"
	"testing"

	interval "github.com/marshauf/coding_challenges/interval_merge/go"
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
		}, {
			name:     "rollup slice bounds out of range",
			input:    []interval.Interval{interval.New(0, 1), interval.New(2, 3), interval.New(3, 4), interval.New(4, 5), interval.New(0, 5)},
			expected: []interval.Interval{interval.New(0, 5)},
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("Merge with %s", row.name), func(t *testing.T) {
			res := interval.Merge(row.input)

			if !equal(res, row.expected) {
				t.Errorf("Expected %v to be merged to %v, got %v", row.input, row.expected, res)
			}
		})
		t.Run(fmt.Sprintf("CMerge with %s", row.name), func(t *testing.T) {
			c := make(chan interval.Interval, 10)
			go (func() {
				for _, interval := range row.input {
					c <- interval
				}
				close(c)
			})()
			res := interval.MergeFromChan(c)
			if !equal(res, row.expected) {
				t.Errorf("Expected %v to be merged to %v, got %v", row.input, row.expected, res)
			}
		})
		t.Run(fmt.Sprintf("MergeFleet with %s", row.name), func(t *testing.T) {
			c := make(chan interval.Interval, 10)
			go (func() {
				for _, interval := range row.input {
					c <- interval
				}
				close(c)
			})()
			res := interval.MergeFleet(c, 8)
			if !equal(res, row.expected) {
				t.Errorf("Expected %v to be merged to %v, got %v", row.input, row.expected, res)
			}
		})

	}
}

func generateBenchmarkInput(size int, seed int64) []interval.Interval {
	r := rand.New(rand.NewSource(seed))
	intervalMaxSize := 10
	intervalMaxStart := size * 100

	intervals := make([]interval.Interval, size)
	for i := range intervals {
		start := r.Intn(intervalMaxStart)
		end := start + r.Intn(intervalMaxSize)
		intervals[i] = interval.New(start, end)
	}
	return intervals
}

func BenchmarkMerge(b *testing.B) {
	table := []struct{ size int }{
		{size: 1},
		{size: 10},
		{size: 100},
		{size: 1_000},
		{size: 10_000},
	}
	for _, row := range table {
		b.Run(fmt.Sprintf("size_%d", row.size), func(b *testing.B) {
			intervals := generateBenchmarkInput(row.size, 0)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				interval.Merge(intervals)
			}
		})
	}
}

func BenchmarkMergeFromChan(b *testing.B) {
	table := []struct{ size int }{
		{size: 1},
		{size: 10},
		{size: 100},
		{size: 1_000},
		{size: 10_000},
	}
	for _, row := range table {
		b.Run(fmt.Sprintf("size_%d", row.size), func(b *testing.B) {
			intervals := generateBenchmarkInput(row.size, 0)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				c := make(chan interval.Interval, row.size)
				go (func() {
					for _, interval := range intervals {
						c <- interval
					}
					close(c)
				})()
				b.StartTimer()
				interval.MergeFromChan(c)
			}
		})
	}
}

func BenchmarkMergeFleet(b *testing.B) {
	table := []struct{ size int }{
		{size: 1},
		{size: 10},
		{size: 100},
		{size: 1_000},
		{size: 10_000},
	}
	for _, row := range table {
		b.Run(fmt.Sprintf("size_%d", row.size), func(b *testing.B) {
			intervals := generateBenchmarkInput(row.size, 0)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				c := make(chan interval.Interval, row.size)
				go (func() {
					for _, interval := range intervals {
						c <- interval
					}
					close(c)
				})()
				b.StartTimer()
				interval.MergeFleet(c, 8)
			}
		})
	}
}

func BenchmarkMergeWithHugeRollup(b *testing.B) {
	intervals := make([]interval.Interval, 1_000)
	for i := range intervals {
		intervals[i] = interval.New(i, i+1)
	}
	intervals[len(intervals)-1] = interval.New(0, len(intervals))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		interval.Merge(intervals)
	}
}
