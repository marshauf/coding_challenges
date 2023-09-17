package interval

// An Interval is a range with a lower and upper bound.
// The lower bound is inclusive, the upper bound is exclusive.
type Interval struct {
	start int // lower bound
	end   int // upper bound
}

// New creates a new Interval with start as lower bound and end as upper bound.
func New(start, end int) Interval {
	return Interval{start, end}
}

// Merges all intervals together which overlap.
// Resulting in a sorted slice of Intervals, where no Interval overlaps another.
// Intervals are bound inclusively below and exclusively above.
// Which means when an Interval upper bound is the same as another Intervals lower bound,
// both intervals don't overlap.
func Merge(intervals []Interval) []Interval {
	merged_intervals := make([]Interval, 0, 8)
	for _, interval := range intervals {
		merged_intervals = merge_into(interval, merged_intervals)
	}
	return merged_intervals
}

// merge_into merges one interval into the provided slice and returns it.
// It either gets prepended if it is smaller than the first element,
// or append if it is larger than the last element.
// If the interval overlaps with another, source is merged into the matched element.
// If the upper bound grew after merge, a rollup happens. See rollup function.
// The length of the slice might grow by one or shrink down to one.
func merge_into(source Interval, intervals []Interval) []Interval {
	for i, interval := range intervals {
		// Does source end before interval?
		if source.end <= interval.start {
			intervals = append([]Interval{source}, intervals...)
			return intervals
		}
		// Is source start inside interval?
		if source.start > interval.start && source.start < interval.end {
			intervals[i].end = max(interval.end, source.end)
			return rollup(intervals, i, intervals[i].end)
		}
		// Is source end inside interval?
		if source.end < interval.end {
			intervals[i].start = min(interval.start, source.start)
		}
		// Is interval inside source?
		if interval.start >= source.start && interval.start < source.end {
			intervals[i] = source
			return rollup(intervals, i, intervals[i].end)
		}
	}
	// source is larger than any interval in intervals
	intervals = append(intervals, source)
	return intervals
}

// rollup merges all overlapping intervals together starting at from, which are smaller than end.
// End changes on merge to the higher value of either end or the upper bound of the merged interval.
// from is the index of the interval in intervals, which started the rollup because its upper bound grew.
func rollup(intervals []Interval, from int, end int) []Interval {
	n := -1
	for i, interval := range intervals[from+1:] {
		// is interval smaller than the upper bound?
		if interval.start < end {
			end = max(end, interval.end)
			n = from + 1 + i
		} else {
			// All following intervals are larger than end, break
			break
		}
	}
	// Slice if merge happened
	if n > -1 {
		intervals = append(intervals[:from+1], intervals[n+1:]...)
	}
	intervals[from].end = end
	return intervals
}

// min returns the smallest number of n or m
func min(n, m int) int {
	if n < m {
		return n
	}
	return m
}

// max returns the largest number of n or m
func max(n, m int) int {
	if n > m {
		return n
	}
	return m
}
