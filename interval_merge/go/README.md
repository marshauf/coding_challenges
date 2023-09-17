# Interval merge in Go

Implementation of the merge interval challenge with Go. The function is called `Merge`.

The functions `MergeStream` and `MergeFleet` are experiments with using channels as input and
goroutines as workers. Based on benchmark feedback, they perform under most scenarios significant worse than the `Merge` function.
The `MergeFleet` also has a bug which randomly appears. Seems like a race condition.

```
    --- FAIL: TestMerge/MergeFleet_with_negative_ranges (0.00s)
        merge_test.go:93: Expected [{-4 0} {-1 2}] to be merged to [{-4 2}], got [{-4 0}]
```

## Requirements

[Go](https://go.dev/doc/install)

## Run

Only tests and benchmarks available. Run the tests with

```bash
go test
```

and the benchmarks with 
```bash
go test -bench Benchmark*
```

