
/*
* Merge intervals together.
* intervals should be a list of intervals.
* Each interval is represented by an Array of length 2.
* The first element is the lower bound and the second element is the upper bound.
* Operates on the provided arguments and modifies it. The modified list is returned.
*/
function merge(intervals) {
    for (let i = 0; i < intervals.length; i++) {
        let interval = intervals[i];
        for (let j = i + 1; j < intervals.length;) {
            if (interval[0] <= intervals[j][0] && interval[1] >= intervals[j][0]) {
                interval[1] = Math.max(interval[1], intervals[j][1]);
                intervals.splice(j, 1);
            } else {
                j++;
            }
        }
    }
    return intervals;
}

function test() {
    const assert = require('node:assert').strict;

    console.log("Example solution");
    let input = [[25, 30], [2, 19], [14, 23], [4, 8]];
    let expected = [[25, 30], [2, 23]];
    let res = merge(structuredClone(input));


    console.log("input", input);
    console.log("expected", expected);
    console.log("got", res);

    assert.deepEqual(res, expected);

    console.log("This implementation does not merge intervals until no interval overlaps another");
    console.log(merge([[25, 30], [2, 5], [7, 23], [4, 8]]));
}

test();
