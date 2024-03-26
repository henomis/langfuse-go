package observer

import "sync"

type queue[T any] struct {
	sync.Mutex
	items []T
}

func (q *queue[T]) Enqueue(item T) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, item)
}

func (q *queue[T]) Dequeue() T {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *queue[T]) Len() int {
	q.Lock()
	defer q.Unlock()
	return len(q.items)
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{}
}

func (q *queue[T]) Clear() {
	q.Lock()
	defer q.Unlock()
	q.items = []T{}
}

func (q *queue[T]) All() []T {
	q.Lock()
	defer q.Unlock()
	items := q.items
	q.items = []T{}
	return items
}
