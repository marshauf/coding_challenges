package interval

type Interval struct {
	start int
	end   int
}

func New(start, end int) Interval {
	return Interval{start, end}
}

func (i Interval) contains(point int) bool {
	return i.start < point && point < i.end
}

func Merge(intervals []Interval) []Interval {
	merged_intervals := make([]Interval, 0, 8)
	for _, interval := range intervals {
		merged_intervals = merge_into(interval, merged_intervals)
	}
	return merged_intervals
}

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
	intervals = append(intervals, source)
	return intervals
}

func rollup(intervals []Interval, from int, end int) []Interval {
	for i, interval := range intervals[from+1:] {
		if interval.start < end {
			end = max(end, interval.end)
			// Re-slicing is expensive, could be optimized by storing the to be deleted ranges
			intervals = append(intervals[:i+from+1], intervals[i+from+2:]...)
		}
	}
	intervals[from].end = end
	return intervals
}

func min(n, m int) int {
	if n < m {
		return n
	}
	return m
}

func max(n, m int) int {
	if n > m {
		return n
	}
	return m
}
