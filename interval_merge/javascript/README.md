# Interval merge Spike in JavaScript

Exploration of the merge interval challenge with JavaScript.
It only solves the example input and doesn't implement all requirements and edge cases.

## Requirements

Requires Node v17 or above to execute because of the use of structuredClone.

## Run

Run by calling merge.js with node. No input required.

```bash
node merge.js
```

Prints out
```bash
Example solution
input [ [ 25, 30 ], [ 2, 19 ], [ 14, 23 ], [ 4, 8 ] ]
expected [ [ 25, 30 ], [ 2, 23 ] ]
got [ [ 25, 30 ], [ 2, 23 ] ]
This implementation does not merge intervals until no interval overlaps another
[ [ 25, 30 ], [ 2, 8 ], [ 7, 23 ] ]
```

## Exploration result

## General

This implementation modifies the input. Depending on the implementation language a detailed documentation is necessary to warn about the modification.
It also requires the complete input beeing in memory all at once. Does not use concurrency.
Depending on the runtime and settings this implementation will not work with inputs in the giga byte range.
The generated garbage should be relative small and not the limiting factor, since it operates on the input. This is untested and an assumption.

### Example input

The example input can be merged completly with one iteration. The first element doesn't get modified.
The second element gets merged with all the following elements. Resulting in a correct result.

### Modified example input

However with an input of `[[25, 30], [2, 5], [7, 23], [4, 8]]` the result is `[ [ 25, 30 ], [ 2, 8 ], [ 7, 23 ] ]`.
Clearly the second and third element overlap. The result is wrong. 

With only one iteration `[2,5]` is merged with `[4,8]` resulting in `[2,8]`. The state after this merge is `[[25,30],[2,8],[7,23]]`.
The merge iterations next and last element beeing `[7,23]`, which has no following element. Resulting in not beeing merged with `[2,8]`.

Possible solutions/thoughts:

* Repeat the merge iteration until the list doesn't change
    * Quite a brute force solution, resulting in a lot of iterations over the input and one last iteration
    which doesn't even modify the list.
    * Would keep the memory consumption low
* Use a second sorted result list of intervals
    * Sorting the result list allows early iteration breaks on inserting or merging iterations, keeps number and length of iterations low
    * Requires more memory than iterating on the input, but it also doesn't modify the input
    * The result will be sorted from small to largest upper bound, not in the order of the input.
    The order isn't listed as a requirement, but should be verified if it is a requirement.
    I assume it isn't a requirement, might also be nice to get it sorted.
* Instead of using a sorted list a BTree might even be faster, depends of course of the implementation of the list/vector or BTree
* The input could be chunked and processed concurrently. At the end merged together.
    * For larger inputs and an implementation which allows the usage of multiple cores this might be faster
* Instead of receiving the input as a list, receive intervals via an iterator/stream/channel
    * Would reduce memory consumption and would allow processing of inputs larger than the available memory
    * Can be combined with chunked processing
    * If receiving data as an iterator is an uncommon practice, the different function interface should be agreed upon
