# wrr

## What
the weight-round-robin strategy by GO
## Features

- Four built-in priorities and supports custom priorities.   
```markdown
LowPriority: with weight = 20

MediumPriority: with weight = 30

HighPriority: with weight = 50

ForcePriority: always first priority
```
```markdown
YourCustomPriority := &Priority{"custom", false, 20}

YourCustomForcePriority := &Priority{"customForce", true, 0}
```
## How
```go
// your node struct
type node struct {
    key   string
    value string
}
```
```go
// New Wrr with Priorities
wrr, err := New[node]().AddPriorities(LowPriority, MediumPriority, HighPriority).Build()
if err != nil {
    panic(err)
}

// Force priority
wrr.AddPriority(ForcePriority)
```
```go
// Push your node with Priority
pe := wrr.Push(LowPriority, &node{
    key:   string(rune(i)),
    value: "node of " + LowPriority.Name,
})
if pe != nil {
    panic(pe)  
}

// Force priority
wrr.Push(ForcePriority, &node{
    key:   string(rune(i)),
    value: "node of " + ForcePriority.Name,
})
```
```go
// Pop your node by weight-round-robin strategy
obj, err := wrr.Pop()
if errors.Is(err, ErrEmptyData) {
    next = false
} else {
	handleYourLogic(obj)
}
```
## Benchmark
### BenchmarkPush  
BenchmarkPush-8   	19523636	        54.03 ns/op	      81 B/op	       1 allocs/op  
### BenchmarkPop  
BenchmarkPop-8   	18534308	        64.47 ns/op	       0 B/op	       0 allocs/op