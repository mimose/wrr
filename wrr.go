package wrr

import "sync"

type Wrr[T any] struct {
	queues      map[string]*queue[T]
	forceQueues *queue[T]
	totalWeight int
	lock        *sync.Mutex
}

type Builder[T any] struct {
	wrr *Wrr[T]
}

func New[T any]() *Builder[T] {
	return &Builder[T]{
		wrr: &Wrr[T]{
			queues:      make(map[string]*queue[T]),
			forceQueues: nil,
			totalWeight: 0,
			lock:        &sync.Mutex{},
		},
	}
}

func (builder *Builder[T]) AddPriority(priority *Priority) *Builder[T] {
	builder.wrr.queues[priority.Name] = &queue[T]{
		weight:     priority.Weight,
		currWeight: 0,
		lock:       &sync.Mutex{},
		objs:       make([]*T, 0),
	}
	builder.wrr.totalWeight += priority.Weight

	return builder
}

func (builder *Builder[T]) AddPriorities(priorities ...*Priority) *Builder[T] {
	for _, priority := range priorities {
		builder.AddPriority(priority)
	}

	return builder
}

func (builder *Builder[T]) Build() (*Wrr[T], error) {
	if len(builder.wrr.queues) == 0 {
		return nil, ErrUncompleted
	}
	return builder.wrr, nil
}

func (wrr *Wrr[T]) Push(p *Priority, obj *T) error {
	if p == nil || obj == nil {
		return ErrIllegalParams
	}

	var q *queue[T]
	if p.Force {
		q = wrr.forceQueues
	} else {
		q = wrr.queues[p.Name]
	}

	if q == nil {
		return ErrUncompleted
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	q.objs = append(q.objs, obj)

	return nil
}

func (wrr *Wrr[T]) Pop() (*T, error) {
	var q *queue[T]
	if wrr.forceQueues != nil && len(wrr.forceQueues.objs) > 0 {
		q = wrr.forceQueues
		return q.pop(nil)
	}

	wrr.lock.Lock()
	defer wrr.lock.Unlock()

	for _, qs := range wrr.queues {
		if len(qs.objs) == 0 {
			continue
		}
		qs.currWeight += qs.weight
		if q == nil || q.currWeight < qs.currWeight {
			q = qs
		}
	}

	if q == nil {
		return nil, ErrEmptyData
	}

	return q.pop(func() {
		q.currWeight -= wrr.totalWeight
	})
}
