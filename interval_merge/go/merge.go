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
	for _, interval := range intervals {
		intervals = merge_into(interval, intervals)
	}
	return intervals
}

func merge_into(source Interval, intervals []Interval) []Interval {

	for _, interval := range intervals {
		if source.end < interval.start {
			intervals = append([]Interval{source}, intervals...)
		}
	}
	return intervals
}
