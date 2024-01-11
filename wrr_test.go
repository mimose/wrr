package wrr

import (
	"errors"
	"sync"
	"testing"
)

type node struct {
	key   string
	value string
}

func Test(t *testing.T) {
	wrr, err := New[node]().AddPriorities(LowPriority, MediumPriority, HighPriority).Build()

	if err != nil {
		t.Fatalf("error: %v", err)
	}

	// push
	pushBasic(100, wrr)

	//next := true
	//for next {
	//	obj, err := wrr.Pop()
	//	if errors.Is(err, ErrEmptyData) {
	//		next = false
	//	} else {
	//		t.Logf("obj: %s : %s, error: %+v", obj.key, obj.value, err)
	//	}
	//}

	group := sync.WaitGroup{}
	group.Add(310)
	for i := 0; i < 310; i++ {
		go func() {
			obj, err := wrr.Pop()
			if errors.Is(err, ErrEmptyData) {
				t.Log("empty")
			} else {
				t.Logf("obj: %s : %s, error: %+v", obj.key, obj.value, err)
			}
			group.Done()
		}()
	}
	group.Wait()

	t.Logf("done")
}

func pushBasic(perSize int, wrr *Wrr[node]) {
	group := sync.WaitGroup{}
	group.Add(3)
	go func() {
		for i := 0; i < perSize; i++ {
			_ = wrr.Push(LowPriority, &node{
				key:   LowPriority.Name,
				value: LowPriority.Name,
			})
		}
		group.Done()
	}()
	go func() {
		for i := 0; i < perSize; i++ {
			_ = wrr.Push(MediumPriority, &node{
				key:   MediumPriority.Name,
				value: MediumPriority.Name,
			})
		}
		group.Done()
	}()
	go func() {
		for i := 0; i < perSize; i++ {
			_ = wrr.Push(HighPriority, &node{
				key:   HighPriority.Name,
				value: HighPriority.Name,
			})
		}
		group.Done()
	}()

	group.Wait()
}

func BenchmarkPush(b *testing.B) {
	wrr, err := New[node]().AddPriorities(LowPriority, MediumPriority, HighPriority).Build()

	if err != nil {
		b.Fatalf("error: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pe := wrr.Push(LowPriority, &node{
			key:   LowPriority.Name,
			value: LowPriority.Name,
		})
		if pe != nil {
			b.Errorf("error: %v", pe)
		}
	}

	b.ReportAllocs()
}

func BenchmarkPop(b *testing.B) {
	wrr, err := New[node]().AddPriorities(LowPriority, MediumPriority, HighPriority).Build()

	if err != nil {
		b.Fatalf("error: %v", err)
	}

	pushBasic(30000000, wrr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = wrr.Pop()
		if errors.Is(err, ErrEmptyData) {
			b.Log("empty")
		}
	}

	b.ReportAllocs()
}
