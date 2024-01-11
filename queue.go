package wrr

import "sync"

type queue[T any] struct {
	weight     int
	currWeight int
	lock       *sync.Mutex
	objs       []*T
}

func (q *queue[T]) pop(fc func()) (*T, error) {
	var obj *T

	if len(q.objs) == 0 {
		return obj, ErrEmptyData
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	obj = q.objs[0]
	q.objs = q.objs[1:]

	fc()

	return obj, nil
}
