use std::ops::Range;

// Merges all intervals together, starting at from, which are inside bound.
// Expects intervals to be sorted.
fn rollup<T>(mut intervals: Vec<Range<T>>, from: usize, mut bound: Range<T>) -> Vec<Range<T>>
where
    T: Ord + Copy,
{
    // Starting merging intervals into from by filtering out all intervals inside bound and at the
    // end update interval bound at index from with highest merged bound.
    let right = &intervals[from + 1..];
    // Filter out all intervals which are contained inside bound.
    // on match extend upper bound to highest upper bound
    let mut right: Vec<Range<T>> = right
        .iter()
        .filter(|i| {
            if bound.contains(&i.start) || bound.contains(&i.end) {
                bound.end = T::max(bound.end, i.end);
                false
            } else {
                true
            }
        })
        .cloned()
        .collect();
    // Update upper bound of interval, where merging started.
    if let Some(last) = intervals.get_mut(from) {
        last.end = bound.end;
    }
    // NOTE instead of updating intervals Vec, a new vector might be faster
    // Drop old interval data
    intervals.truncate(from + 1);
    // Append new interval data
    intervals.append(&mut right);
    intervals
}

/// Merges all ranges together.
/// Resulting in a sorted Vector, where no Range overlaps another.
/// Ranges are bound inclusively below and exclusively above.
///
/// ````
/// use rs::merge;
///
/// let res = merge(vec![0..4, 2..6, 0..2]);
/// assert_eq!(res, vec![0..6]);
/// ````
pub fn merge<I, T>(intervals: I) -> Vec<Range<T>>
where
    I: IntoIterator<Item = Range<T>>,
    T: Ord + Copy, // Order is required for comparing bounds, Copy is required to generate a new
                   // Vec<Range<T>>.
{
    intervals
        .into_iter()
        .fold(Vec::new(), |mut acc: Vec<Range<T>>, interval| {
            for (n, i) in acc.iter_mut().enumerate() {
                if interval.end <= i.start {
                    // interval is below current and all following intervals
                    // insert interval at this position and shift following intervals to the right
                    acc.insert(n, interval);
                    return acc;
                }
                if i.contains(&interval.start) {
                    // interval starts inside current interval
                    // extend upper bound of current interval to include range of interval
                    i.end = T::max(i.end, interval.end);
                    // merge all ranges which are inside the bound
                    let bound = i.clone();
                    return rollup(acc, n, bound);
                }
                if i.contains(&interval.end) {
                    // interval ends inside current interval, but was not part of any previous one.
                    // extend lower bound to include interval range
                    i.start = T::min(i.start, interval.start);
                    // rollup is not necessary, since interval upper bound ends inside current
                    // interval
                    return acc;
                }
                if interval.contains(&i.start) {
                    // interval is bigger than current interval
                    // extend current interval ranges to bigger interval
                    i.start = interval.start;
                    i.end = interval.end;
                    // since upper bound of interval is not inside current interval rollup
                    let bound = i.clone();
                    return rollup(acc, n, bound);
                }
            }
            // first iteration of fold or intervals lower bound is bigger than any folded upper
            // interval bound
            acc.push(interval);
            acc
        })
}

#[cfg(test)]
mod tests {
    use crate::merge;

    #[allow(clippy::single_range_in_vec_init)]
    #[test]
    fn test_integer_ranges() {
        let table = vec![
            // empty
            (vec![], vec![]),
            // one
            (vec![0..1], vec![0..1]),
            // Ranges which do not overlap
            (vec![1..2, 3..6, 8..12], vec![1..2, 3..6, 8..12]),
            // example
            (vec![25..30, 2..19, 14..23, 4..8], vec![2..23, 25..30]),
            // example with a range at the end which includes all ranges before it
            (vec![25..30, 2..19, 14..23, 4..8, 0..40], vec![0..40]),
            // Range at the end merges multiple previous ranges
            (vec![0..2, 20..22, 3..5, 7..8, 1..8], vec![0..8, 20..22]),
            // negative ranges
            (vec![-4..0, -1..2], vec![-4..2]),
        ];
        for row in table {
            let res = merge(row.0.clone());
            assert_eq!(
                res, row.1,
                "expected merge {:?} to be {:?} but got {:?}",
                row.0, row.1, res
            );
        }
    }

    #[allow(clippy::single_range_in_vec_init)]
    #[test]
    fn test_char_ranges() {
        let table = vec![
            // empty
            (vec![], vec![]),
            // one
            (vec!['a'..'b'], vec!['a'..'b']),
            // Ranges which do not overlap
            (
                vec!['a'..'b', 'c'..'d', 'f'..'g'],
                vec!['a'..'b', 'c'..'d', 'f'..'g'],
            ),
            // Overlapping
            (vec!['a'..'d', 'c'..'d', 'a'..'g'], vec!['a'..'g']),
        ];
        for row in table {
            let res = merge(row.0.clone());
            assert_eq!(
                res, row.1,
                "expected merge {:?} to be {:?} but got {:?}",
                row.0, row.1, res
            );
        }
    }
}
